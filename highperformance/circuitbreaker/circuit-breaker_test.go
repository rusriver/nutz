package circuitbreaker_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/rusriver/nutz/highperformance/circuitbreaker"
	loggeterrzl "github.com/rusriver/nutz/loggeterr-zl"
	"github.com/stretchr/testify/require"
)

// this test takes 3.5 minutes,
// see https://stackoverflow.com/questions/73033725/how-to-edit-default-go-test-timeout-in-visual-studio-code
func Test_001(t *testing.T) {
	tracer := loggeterrzl.New(nil)
	ctx, ctxCancel := context.WithCancel(context.Background())

	cb_opts_defaults := func(cb *circuitbreaker.CircuitBreaker) {
		cb.ChanReportSize = 32
		cb.ChanReigniteSize = 2
		cb.ReignitePeriod = time.Second * 10
		cb.SwitchOnWhenAbove = 3 // intentionally bigger
		cb.SwitchOffWhenBelow = -3
		cb.PositiveAutoPullUpPerSec = 0.1
		cb.OffWaitTime = time.Minute * 1
		cb.Context = ctx
	}

	cb := circuitbreaker.New(cb_opts_defaults, func(cb *circuitbreaker.CircuitBreaker) {
		cb.DebugTracer = tracer.NSTIC("CB-01").ITracer()
	})

	tracer = tracer.NSTIC("main").ITracer()

	chDQF := make(chan func(), 32)
	tkr01 := time.NewTicker(time.Millisecond * 500)

	G1_Results := true

	var (
		stats_executed = 0
		stats_missed   = 0
	)

	go func() {
		for {
			select {
			case <-ctx.Done():
				return

			case f := <-chDQF:
				f()
				for i := 0; i < len(chDQF); i++ {
					f = <-chDQF
					f()
				}

			case <-tkr01.C:
				_, execd := cb.WantToDo(func() (success bool) {
					tracer.InfoEvent("45").SendMsg("tick G1", G1_Results)
					stats_executed++
					return G1_Results
				})
				if !execd {
					tracer.InfoEvent("50").SendMsg("tick void")
					stats_missed++
				}

			} // select
		} // for
	}()

	time.Sleep(time.Second * 10)

	tracer.InfoEvent("59").SendMsg("+++ 10 sec")

	time.Sleep(time.Second * 5)

	tracer.InfoEvent("63").SendMsg("+++ 5 sec")

	chDQF <- func() {
		G1_Results = false
		tracer.InfoEvent("67").SendMsg("+++ G1_Results", G1_Results)
	}

	time.Sleep(time.Second * 30)

	chDQF <- func() {
		G1_Results = true
		tracer.InfoEvent("74").SendMsg("+++ G1_Results", G1_Results)
	}

	time.Sleep(time.Second * 30)

	chDQF <- func() {
		G1_Results = false
		tracer.InfoEvent("81").SendMsg("+++ G1_Results", G1_Results)
	}

	time.Sleep(time.Second * 30)

	chDQF <- func() {
		G1_Results = true
		tracer.InfoEvent("88").SendMsg("+++ G1_Results", G1_Results)
	}

	time.Sleep(time.Second * 90)

	ctxCancel()
	fmt.Println("stats_executed", stats_executed)
	fmt.Println("stats_missed", stats_missed)

	require.InDelta(t, 55, stats_executed, 5)
	require.InDelta(t, 335, stats_missed, 5)
}
