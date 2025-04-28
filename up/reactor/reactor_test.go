package reactor_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/lithammer/shortuuid/v4"
	"github.com/rusriver/nutz/controlflow"
	"github.com/rusriver/nutz/logger"
	"github.com/rusriver/nutz/up/reactor"
	"github.com/stretchr/testify/require"
)

const (
	Subsystem01 = "SUBSYSTEM-01"
	Subsystem02 = "SUBSYSTEM-02"
	Subsystem03 = "SUBSYSTEM-03"
)

func Test_Reactor_01(t *testing.T) {
	logger := logger.New(shortuuid.New()[:8], nil)

	r := &reactor.Reactor{
		Name:    "R1",
		Context: context.Background(),
		F: func(r *reactor.Reactor, m *reactor.Msg) {
			switch m.Signal {
			case reactor.Signal_End:
				r.Logger.InfoEvent("22").SendMsg("end")

			case reactor.Signal_Ready, reactor.Signal_Cede:
				r.Logger.InfoEvent("24").SendMsg("ready signal")
				status := r.Status()
				r.Logger.InfoEvent("29").
					Strs("fresh", status.Fresh).
					Strs("ready", status.Ready).
					Strs("ceded", status.Ceded).
					SendMsg(m.PId + " signalled")
				if status.AllDone {
					r.Logger.InfoEvent("35").SendMsg("DO THE NEXT STEP")
					r.End()
					// do something
				}
			}
		},
	}
	r.Logger = logger.NSLIC("1e0194").Str("name", r.Name).ILogger()
	reactor.Start(r, []string{Subsystem01, Subsystem02, Subsystem03})

	fmt.Println("STARTED")

	time.Sleep(time.Second * 1)

	r.Send(&reactor.Msg{reactor.Signal_Ready, Subsystem01})

	time.Sleep(time.Second * 1)

	r.Send(&reactor.Msg{reactor.Signal_Cede, Subsystem02})

	time.Sleep(time.Second * 1)

	r.Send(&reactor.Msg{reactor.Signal_Ready, Subsystem03})

	time.Sleep(time.Second * 1)
	fmt.Println("COMPLETE")
}

// lightweight version
func Test_Reactor_02(t *testing.T) {
	logger := logger.New(shortuuid.New()[:8], nil)

	r := &reactor.Reactor{
		Name:    "R1",
		Context: context.Background(),
		F: func(r *reactor.Reactor, m *reactor.Msg) {
			switch m.Signal {
			case reactor.Signal_Ready:

				r.Logger.InfoEvent("24").SendMsg("ready signal")
				status := r.Status()
				r.Logger.InfoEvent("29").
					Strs("fresh", status.Fresh).
					Strs("ready", status.Ready).
					Strs("ceded", status.Ceded).
					SendMsg(m.PId + " ready")
				if status.AllDone {
					r.Logger.InfoEvent("35").SendMsg("DO THE NEXT STEP")
					r.End()
					// do something
				}

			} //switch
		},
	}
	r.Logger = logger.NSLIC("1e0194").Str("name", r.Name).ILogger()
	reactor.Start(r, []string{Subsystem01, Subsystem02, Subsystem03})

	fmt.Println("STARTED")

	time.Sleep(time.Second * 1)

	r.Send(&reactor.Msg{reactor.Signal_Ready, Subsystem01})

	time.Sleep(time.Second * 1)

	r.Send(&reactor.Msg{reactor.Signal_Cede, Subsystem02})

	time.Sleep(time.Second * 1)

	r.Send(&reactor.Msg{reactor.Signal_Ready, Subsystem03})

	time.Sleep(time.Second * 1)
	fmt.Println("COMPLETE")
}

func Test_Reactor_03(t *testing.T) {
	logger := logger.New(shortuuid.New()[:8], nil)

	controlflow.Try(func() (err error) {

		r := &reactor.Reactor{
			Name:    "R1",
			Context: context.Background(),
			F: func(r *reactor.Reactor, m *reactor.Msg) {
				switch m.Signal {
				case reactor.Signal_End:
					r.Logger.InfoEvent("22").SendMsg("end")

				case reactor.Signal_Ready, reactor.Signal_Cede:
					r.Logger.InfoEvent("24").SendMsg("ready signal")
					status := r.Status()
					r.Logger.InfoEvent("29").
						Strs("fresh", status.Fresh).
						Strs("ready", status.Ready).
						Strs("ceded", status.Ceded).
						SendMsg(m.PId + " signalled")
					if status.AllDone {
						r.Logger.InfoEvent("35").SendMsg("DO THE NEXT STEP")
						r.End()
						// do something
					}
				}
			},
		}
		r.Logger = logger.NSLIC("1e0194").Str("name", r.Name).ILogger()
		reactor.Start(r, []string{Subsystem01, Subsystem02, Subsystem03})

		fmt.Println("STARTED")

		time.Sleep(time.Second * 1)

		r.Send(&reactor.Msg{Signal: reactor.Signal_End})

		time.Sleep(time.Second * 1)

		r.End() // this must be OK

		// panic test
		r.Send(&reactor.Msg{reactor.Signal_Ready, Subsystem01})

		return
	}).Catch(func(e *controlflow.Exception) {
		fmt.Println(e)
		require.Equal(t, "send on closed channel; 01967ad5b672bff68807 reactor is not inited or has been ended", e.String())
	})
}
