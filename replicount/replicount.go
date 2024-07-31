package replicount

import (
	"context"
	"fmt"
	"time"

	"github.com/rusriver/nutz/timmer"
	"github.com/rusriver/ttlcache/v3"
)

type Replicount struct {
	SlowPollPeriodPerReplica time.Duration // basic period, per-replica, in slow mode
	TTLMultiple              int           // multiple of basic poll period; also the TTL of slow mode cache
	FastModeSpeedMultiple    int           // relative to basic poll period
	PollFunc                 PollFunc
	ReportResultsFunc        func(newResult *ChangeableObject)
	Context                  context.Context

	DQF                 chan func(r *Replicount)
	DQF_Len             int
	slowCache           *ttlcache.Cache[string, bool]
	fastCache           *ttlcache.Cache[string, bool]
	slowTTL             time.Duration
	fastTTL             time.Duration
	state               int
	mainTimer           *timmer.Timmer
	fastModeLengthTimer *timmer.Timmer
	fastModeLength      time.Duration
	currentResult       *ChangeableObject
}

type PollFunc func() (idPtr *string)

// must be inited by user, and passed to consumers of this data
type DataContainer struct {
	ChangeableObject *ChangeableObject
}

type ChangeableObject struct {
	NumberOfReplicas int
	ListOfReplicas   []string
}

const (
	mode_Slow = iota
	mode_Fast
)

func New(of ...func(r *Replicount)) (r *Replicount) {
	r = &Replicount{
		SlowPollPeriodPerReplica: time.Second * 30,
		TTLMultiple:              6,
		FastModeSpeedMultiple:    20,
		Context:                  context.Background(),
		DQF_Len:                  24,
		state:                    mode_Slow,
		mainTimer:                timmer.New(),
		fastModeLengthTimer:      timmer.New(),
		currentResult:            &ChangeableObject{NumberOfReplicas: 1},
	}
	for _, f := range of {
		f(r)
	}
	r.DQF = make(chan func(r *Replicount), r.DQF_Len)

	r.slowTTL = r.SlowPollPeriodPerReplica * time.Duration(r.TTLMultiple)
	r.fastModeLength = r.slowTTL
	r.fastTTL = r.slowTTL / time.Duration(r.FastModeSpeedMultiple)

	r.slowCache = ttlcache.New[string, bool](ttlcache.WithDisableTouchOnHit[string, bool]())
	go r.slowCache.Start()

	r.fastCache = ttlcache.New[string, bool](ttlcache.WithDisableTouchOnHit[string, bool]())
	go r.fastCache.Start()

	// the main worker
	go func() {
		for {
			select {
			case <-r.Context.Done():
				r.slowCache.Stop()
				r.fastCache.Stop()
				return
			case dqf := <-r.DQF:
				dqf(r)
			}
		}
	}()

	go r.scheduler()

	return
}

func (r *Replicount) scheduler() {
	r.state = mode_Slow
	r.mainTimer.Restart(time.Second)
	for {
		select {
		case <-r.Context.Done():
			return
		case <-r.mainTimer.C:
			if r.PollFunc != nil {
				switch r.state {
				case mode_Slow:
					go func() {
						idPtr := r.PollFunc()
						r.DQF <- func(r *Replicount) {
							void, newSlow, _ := r.accoutForTheReplica(idPtr)
							if void {
								return
							}

							cnt := r.slowCache.Len()
							if cnt < 1 {
								cnt = 1
							}
							if cnt != r.currentResult.NumberOfReplicas {
								newRes := &ChangeableObject{
									NumberOfReplicas: cnt,
									ListOfReplicas:   r.slowCache.Keys(),
								}
								r.currentResult = newRes
								if r.ReportResultsFunc != nil {
									r.ReportResultsFunc(newRes)
								}
							}

							if newSlow {
								r.state = mode_Fast
								for _, k := range r.slowCache.Keys() {
									r.fastCache.SetWithTTL(k, true, r.fastTTL/time.Duration(r.currentResult.NumberOfReplicas*2))
								}
								r.mainTimer.Restart( // immediate race-restart of the main timer
									r.SlowPollPeriodPerReplica /
										time.Duration(r.currentResult.NumberOfReplicas) /
										time.Duration(r.FastModeSpeedMultiple),
								)
								r.fastModeLengthTimer.Restart(r.fastModeLength)
							}

						} // dqf func
						d := r.SlowPollPeriodPerReplica /
							time.Duration(r.currentResult.NumberOfReplicas)
						r.mainTimer.Restart(d)
					}() // go func

				case mode_Fast:
					go func() {
						idPtr := r.PollFunc()
						r.DQF <- func(r *Replicount) {
							void, _, newFast := r.accoutForTheReplica(idPtr)
							if void {
								return
							}
							if newFast {
								fmt.Println("-- new fast / restart")
								r.fastModeLengthTimer.Restart(r.fastModeLength)
							}

							cnt := r.fastCache.Len()
							if cnt < 1 {
								cnt = 1
							}
							if cnt != r.currentResult.NumberOfReplicas {
								newRes := &ChangeableObject{
									NumberOfReplicas: cnt,
									ListOfReplicas:   r.fastCache.Keys(),
								}
								r.currentResult = newRes
								if r.ReportResultsFunc != nil {
									r.ReportResultsFunc(newRes)
								}
							}
						} // dqf func
						d := r.SlowPollPeriodPerReplica /
							time.Duration(r.currentResult.NumberOfReplicas) /
							time.Duration(r.FastModeSpeedMultiple)
						r.mainTimer.Restart(d)
					}() // go func

				} // switch
			} // if r.PollFunc != nil

		case <-r.fastModeLengthTimer.C:
			r.DQF <- func(r *Replicount) {
				fmt.Println("-- slow mode")
				r.state = mode_Slow
			}

		} // select
	} // for
}

func (r *Replicount) accoutForTheReplica(idPtr *string) (void, newSlow, newFast bool) {
	if idPtr == nil || len(*idPtr) == 0 {
		void = true
		return
	}

	item := r.slowCache.Get(*idPtr)
	if item == nil {
		newSlow = true
	}
	r.slowCache.SetWithTTL(*idPtr, true, r.slowTTL)

	item = r.fastCache.Get(*idPtr)
	if item == nil {
		newFast = true
	}
	r.fastCache.SetWithTTL(*idPtr, true, r.fastTTL)

	return
}
