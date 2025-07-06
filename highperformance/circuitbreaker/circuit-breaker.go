package circuitbreaker

import (
	"context"
	"time"

	"github.com/rusriver/nutz/controlflow"
	loggeterrapi "github.com/rusriver/nutz/loggeterr-api"
	"github.com/rusriver/nutz/loggeterr-api/trace"
)

const (
	// don't use 0, because 0 is reserved as an indication of missing data in metrics
	CbState_OFF_DENIED = 1
	CbState_ON_ALLOWED = 2
	//
	// don't use 0, because 0 is reserved as an indication of missing data in metrics
	GState_DENY_RESTART = 1_00
	GState_ALLOW        = 2_00
	GState_DENY_WAIT    = 3_00
)

type CircuitBreaker struct {
	ChanReportSize           int                  // reporting buffer size; note, it doesn't affect the algorithm
	ChanReigniteSize         int                  // N of reignition tokens written every ReignitePeriod; typically is <= SwitchOnWhenAbove
	ReignitePeriod           time.Duration        // period to refill the chReignite, when in DENY_RESTART mode
	SwitchOnWhenAbove        float32              //
	SwitchOffWhenBelow       float32              // typically should be >= ChanReigniteSize
	PositiveAutoPullUpPerSec float32              // amount of auto-pull-up per second, when in ALLOW mode, from negative zone
	OffWaitTime              time.Duration        // wait time in DENY_WAIT mode
	Context                  context.Context      // can be used to stop it
	DQF_Size                 int                  //
	DebugTracer              loggeterrapi.ITracer // normally nil, set it to enable debug mode
	StateHookFunc            StateHookFunc        // dep injection // this can be used to inject code that does metrics et al
	//
	cbState    int
	chReport   chan bool
	chReignite chan struct{}
	chDQF      chan func()
}

type StateHookFunc func(cbState, gState int)

func New(fOpts ...func(cb *CircuitBreaker)) (cb *CircuitBreaker) {
	cb = &CircuitBreaker{}
	// the defaults, intentionally in assignment form for easier copy-paste
	cb.ChanReportSize = 32
	cb.ChanReigniteSize = 3
	cb.ReignitePeriod = time.Second * 5
	cb.SwitchOnWhenAbove = 3
	cb.SwitchOffWhenBelow = -3
	cb.PositiveAutoPullUpPerSec = 0.1
	cb.OffWaitTime = time.Minute * 1
	cb.Context = context.Background()
	cb.DQF_Size = 32

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
	var (
		counter  float32
		lastTime time.Time
		gState   int
	)
	if cb.DebugTracer != nil {
		cb.DebugTracer.DiagnosticEvent("45").SendMsg("g started")
	}

	// go func() {
	// 	for {
	// 		time.Sleep(500 * time.Millisecond)
	// 		cb.DebugTracer.DiagnosticEvent("82").SendMsg("-- counter", counter)
	// 	}
	// }()

	var setState func(int)
	setState = func(newState int) {
		switch newState {

		case GState_DENY_RESTART:
			gState = GState_DENY_RESTART
			cb.cbState = CbState_OFF_DENIED
			counter = 0
			lastTime = time.Now()
			if cb.DebugTracer != nil {
				cb.DebugTracer.DiagnosticEvent("96").SendMsg("DENY-RESTART")
			}

		case GState_ALLOW:
			gState = GState_ALLOW
			cb.cbState = CbState_ON_ALLOWED
			counter = 0
			lastTime = time.Now()
			if cb.DebugTracer != nil {
				cb.DebugTracer.DiagnosticEvent("105").SendMsg("ALLOW")
			}

		case GState_DENY_WAIT:
			gState = GState_DENY_WAIT
			cb.cbState = CbState_OFF_DENIED
			counter = 0
			lastTime = time.Now()
			go func() {
				select {
				case <-cb.Context.Done():
					return
				case <-time.After(cb.OffWaitTime):
					select {
					case cb.chDQF <- func() {
						setState(GState_DENY_RESTART)
					}:
					default:
						controlflow.RaiseCause(nil, &trace.BacktracePoint{"d64bc28d72f7", "DQF overflow"})
					}
				}
			}()
			if cb.DebugTracer != nil {
				cb.DebugTracer.DiagnosticEvent("130").SendMsg("DENY-WAIT")
			}
		}

		if cb.StateHookFunc != nil {
			cb.StateHookFunc(cb.cbState, gState)
		}
	}

	setState(GState_DENY_RESTART)

	// initial ignition, once
	if cb.DebugTracer != nil {
		cb.DebugTracer.DiagnosticEvent("143").SendMsg("ignition")
	}
	for i := 0; i < cb.ChanReigniteSize; i++ {
		select {
		case cb.chReignite <- struct{}{}:
		default:
		}
	}

	autoPullUpPositive := func() {
		tNow := time.Now()
		if counter < 0 {
			seconds := tNow.Sub(lastTime).Seconds()
			counter += float32(seconds) * cb.PositiveAutoPullUpPerSec
			if counter > 0 {
				counter = 0
			}
		}
		lastTime = tNow
	}

	reigniteTicker := time.NewTicker(cb.ReignitePeriod)

	// main loop here
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
			switch gState {

			case GState_DENY_RESTART:
				autoPullUpPositive() // here for symmetricity
				switch r {
				case true:
					counter++
					if counter >= cb.SwitchOnWhenAbove {
						setState(GState_ALLOW)
					}
				case false:
					counter--
					if counter <= cb.SwitchOffWhenBelow {
						setState(GState_DENY_WAIT)
					}
				}

			case GState_ALLOW:
				autoPullUpPositive()
				switch r {
				case true:
					counter++ // this will not happen, usually, because un-reported
					if counter > 0 {
						counter = 0
					}
				case false:
					counter--
					if counter < cb.SwitchOffWhenBelow {
						setState(GState_DENY_WAIT)
					}
				}

			case GState_DENY_WAIT:
				// ignore any reports, just wait
			}

		case <-reigniteTicker.C: // not DQF, intentionally randomized, relative to the request time
			if gState == GState_DENY_RESTART {
				for i := 0; i < cb.ChanReigniteSize; i++ {
					select {
					case cb.chReignite <- struct{}{}:
					default:
					}
				}
				if cb.DebugTracer != nil {
					cb.DebugTracer.DiagnosticEvent("225").SendMsg("reignition")
				}
			}

		} // select
	} // for
}

func (cb *CircuitBreaker) Report(r bool) {
	switch cb.cbState {
	case CbState_OFF_DENIED:
		select {
		case cb.chReport <- r:
		default:
		}
	case CbState_ON_ALLOWED:
		if !r {
			select {
			case cb.chReport <- r:
			default:
			}
		}
	}
}

func (cb *CircuitBreaker) IsOn() (ok bool) {
	if cb.cbState == CbState_ON_ALLOWED {
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

func (cb *CircuitBreaker) WantToDo(f func() (success bool)) (result bool, executed bool) {
	if cb.IsOn() {
		result = f()
		cb.Report(result)
		return result, true
	} else {
		return false, false
	}
}
