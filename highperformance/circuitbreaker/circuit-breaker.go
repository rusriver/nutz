package circuitbreaker

import (
	"context"
	"errors"
	"log"
	"math/rand"
	"time"
)

const (
	cbState_OpenDenied    = 0
	cbState_ClosedAllowed = 1
	//
	gState_DenyRestart = 1_00
	gState_Allow       = 2_00
	gState_DenyWait    = 3_00
	//
	debugPrefix = "CIRCUIT-BREAKER: "
)

type CircuitBreaker struct {
	ChanReportSize      int
	ChanReigniteSize    int
	ReignitePeriod      time.Duration
	CloseAllowWhenAbove int
	OpenDenyWhenBelow   int
	DenyWaitTime        time.Duration
	Context             context.Context
	DQF_Size            int
	DebugMode           bool
	//
	state      int
	chReport   chan bool
	chReignite chan struct{}
	chDQF      chan func()
}

func New(fOpts ...func(cb *CircuitBreaker)) (cb *CircuitBreaker) {
	cb = &CircuitBreaker{}
	// defaults here, intentionally in assignment form for easier copy-paste
	cb.ChanReportSize = 32
	cb.ChanReigniteSize = 3
	cb.ReignitePeriod = time.Second * 5
	cb.CloseAllowWhenAbove = 3
	cb.OpenDenyWhenBelow = -3
	cb.DenyWaitTime = time.Minute * 1
	cb.Context = context.Background()
	cb.DQF_Size = 32
	cb.DebugMode = false

	for _, f := range fOpts {
		f(cb)
	}

	cb.chReport = make(chan bool, cb.ChanReportSize)
	cb.chReignite = make(chan struct{}, cb.ChanReigniteSize)
	cb.chDQF = make(chan func(), cb.DQF_Size)

	go cb.g()

	return cb
}

func (cb *CircuitBreaker) g() {
	counter := 0
	lastTime := time.Now()
	state := gState_DenyRestart
	tkr := time.NewTicker(cb.ReignitePeriod) // keep it so, don't move to DQF -- it's useful, see docs why

	if cb.DebugMode {
		log.Printf(debugPrefix + "g started\n")
	}

	// initial ignition, once
	rand.Seed(time.Now().UnixNano())
	time.Sleep(time.Duration(rand.Intn(int(cb.ReignitePeriod.Milliseconds()))) * time.Millisecond)
	for i := 0; i < cb.ChanReigniteSize; i++ {
		select {
		case cb.chReignite <- struct{}{}:
		default:
		}
	}
	if cb.DebugMode {
		log.Printf(debugPrefix + "ignition\n")
	}

	pullUpPositiveRPS := func() {
		tNow := time.Now()
		if counter < 0 {
			seconds := tNow.Sub(lastTime).Seconds()
			counter += int(seconds)
			if counter > 0 {
				counter = 0
			}
		}
		lastTime = tNow
	}

	var setState func(int)
	setState = func(newState int) {
		switch newState {
		case gState_DenyRestart:
			cb.state = cbState_OpenDenied
			counter = 0
			lastTime = time.Now()
			state = gState_DenyRestart
			if cb.DebugMode {
				log.Printf(debugPrefix + "DENY-RESTART\n")
			}

		case gState_Allow:
			cb.state = cbState_ClosedAllowed
			counter = 0
			lastTime = time.Now()
			state = gState_Allow
			if cb.DebugMode {
				log.Printf(debugPrefix + "ALLOW\n")
			}

		case gState_DenyWait:
			cb.state = cbState_OpenDenied
			counter = 0
			lastTime = time.Now()
			state = gState_DenyWait
			go func() {
				time.Sleep(cb.DenyWaitTime)
				select {
				case cb.chDQF <- func() {
					setState(gState_DenyRestart)
				}:
				default:
					panic(errors.New("d64bc28d72f7 DQF overflow"))
				}
			}()
			if cb.DebugMode {
				log.Printf(debugPrefix + "DENY-WAIT\n")
			}
		}
	}

	for {
		select {
		case <-cb.Context.Done():
			return

		case f := <-cb.chDQF:
			f()
			for i := 0; i < len(cb.chDQF); i++ {
				f = <-cb.chDQF
				f()
			}

		case r := <-cb.chReport:

			switch state {
			case gState_DenyRestart:
				pullUpPositiveRPS() // here for symmetricity
				switch r {
				case true:
					counter++
					if counter > cb.CloseAllowWhenAbove {
						setState(gState_Allow)
					}
				case false:
					counter--
					if counter < cb.OpenDenyWhenBelow {
						setState(gState_DenyWait)
					}
				}

			case gState_Allow:
				pullUpPositiveRPS()
				switch r {
				case true:
					counter++ // this will not happen, usually
					if counter > 0 {
						counter = 0
					}
				case false:
					counter--
					if counter < cb.OpenDenyWhenBelow {
						setState(gState_DenyWait)
					}
				}

			case gState_DenyWait:
				// ignore any reports, just wait

			}

		case <-tkr.C: // not DQF, intentionally randomized
			if state == gState_DenyRestart {
				for i := 0; i < cb.ChanReigniteSize; i++ {
					select {
					case cb.chReignite <- struct{}{}:
					default:
					}
				}
				if cb.DebugMode {
					log.Printf(debugPrefix + "reignition\n")
				}
			}

		} // select
	} // for
}

func (cb *CircuitBreaker) Report(r bool) {
	switch cb.state {
	case cbState_OpenDenied:
		select {
		case cb.chReport <- r:
		default:
		}
	case cbState_ClosedAllowed:
		if !r {
			select {
			case cb.chReport <- r:
			default:
			}
		}
	}
}

func (cb *CircuitBreaker) IsClosedAllowed() (ok bool) {
	if cb.state == cbState_ClosedAllowed {
		return true
	} else {
		select {
		case <-cb.chReignite:
			return true
		default:
		}
	}
	return false
}

func (cb *CircuitBreaker) WantToDo(f func() (success bool)) (r bool) {
	if cb.IsClosedAllowed() {
		r = f()
		cb.Report(r)
	}
	return false
}
