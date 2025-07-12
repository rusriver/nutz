package reactor_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/rusriver/nutz/controlflow"
	loggeterrzl "github.com/rusriver/nutz/loggeterr-zl"
	"github.com/rusriver/nutz/up/reactor"
	"github.com/stretchr/testify/require"
)

const (
	Subsystem01 = "SUBSYSTEM-01"
	Subsystem02 = "SUBSYSTEM-02"
	Subsystem03 = "SUBSYSTEM-03"
)

func Test_Reactor_01(t *testing.T) {
	tracer := loggeterrzl.New(nil).NSTIC("1e0194").ITracer()

	tracer = tracer.NSTIC("xuid here").ITracer()

	success := false

	r := reactor.NewReactor("R1", context.Background(), tracer)
	r.NeedsPattern = map[string]uint8{
		Subsystem01: reactor.RequiredStatus_UpOrCeded,
		Subsystem02: reactor.RequiredStatus_UpOrCeded,
		Subsystem03: reactor.RequiredStatus_UpOrCeded,
	}
	r.Func = func(r *reactor.Reactor, m *reactor.Msg) {
		switch m.Signal {
		case reactor.Signal_End:
			r.Tracer.LogEvent("22").SendMsg("end")

		case reactor.Signal_IsUp, reactor.Signal_Cede:
			r.Tracer.LogEvent("24").SendMsg("up signal")
			status := r.GetStatus_InsideReactorFuncOnly()
			r.Tracer.LogEvent("29").
				Strs("fresh", status.Fresh).
				Strs("ready", status.Up).
				Strs("ceded", status.Ceded).
				SendMsg(m.Name + " signalled")

			// this is common idiom
			if status.CanReact {
				r.Tracer.LogEvent("45").SendMsg("READINESS ACHIEVED - DO THE NEXT STEP")
				r.Done = true
				success = true
				// do something
			} else if status.NoFreshLeft { // optional
				r.Tracer.LogEvent("47").SendMsg("NO FRESH LEFT - DO THE NEXT STEP")
				r.Done = true
				// do something
			}
		}
	}
	r.StartFresh()

	fmt.Println("STARTED")

	time.Sleep(time.Second * 1)
	r.SendStatus_Up(nil, Subsystem01)

	time.Sleep(time.Second * 1)
	r.SendStatus_Up(nil, Subsystem02)

	time.Sleep(time.Second * 1)
	r.SendStatus_Cede(nil, Subsystem03)

	time.Sleep(time.Second * 1)
	r.SendStatus_Up(nil, "sub-4")

	time.Sleep(time.Second * 1)
	fmt.Println("COMPLETE")

	require.Equal(t, true, success)
}

// lightweight version
func Test_Reactor_02(t *testing.T) {
	tracer := loggeterrzl.New(nil).NSTIC("1e0194").ITracer()

	success := false

	r := reactor.NewReactor("R1", context.Background(), tracer)
	r.NeedsPattern = map[string]uint8{
		Subsystem01: reactor.RequiredStatus_UpOrCeded,
		Subsystem02: reactor.RequiredStatus_UpOrCeded,
		Subsystem03: reactor.RequiredStatus_UpOrCeded,
	}
	r.Func = func(r *reactor.Reactor, m *reactor.Msg) {
		switch m.Signal {
		case reactor.Signal_IsUp:

			r.Tracer.LogEvent("24").SendMsg("ready signal")
			status := r.GetStatus_InsideReactorFuncOnly()
			r.Tracer.LogEvent("29").
				Strs("fresh", status.Fresh).
				Strs("ready", status.Up).
				Strs("ceded", status.Ceded).
				SendMsg(m.Name + " ready")

			// this is common idiom
			if status.CanReact {
				r.Tracer.LogEvent("45").SendMsg("READINESS ACHIEVED - DO THE NEXT STEP")
				r.Done = true
				success = true
				// do something
			} else if status.NoFreshLeft {
				r.Tracer.LogEvent("47").SendMsg("NO FRESH LEFT - DO THE NEXT STEP")
				r.Done = true
				// do something
			}
		} //switch
	}
	r.StartFresh()

	fmt.Println("STARTED")

	time.Sleep(time.Second * 1)
	r.SendStatus_Up(nil, Subsystem01)

	time.Sleep(time.Second * 1)
	r.SendStatus_Cede(nil, Subsystem02)

	time.Sleep(time.Second * 1)
	r.SendStatus_Up(nil, Subsystem03)

	time.Sleep(time.Second * 1)
	fmt.Println("COMPLETE")

	require.Equal(t, true, success)
}

func Test_Reactor_03(t *testing.T) {
	tracer := loggeterrzl.New(nil).NSTIC("1e0194").ITracer()

	controlflow.Try(func() (err any) {

		r := reactor.NewReactor("R1", context.Background(), tracer)
		r.NeedsPattern = map[string]uint8{
			Subsystem01: reactor.RequiredStatus_UpOrCeded,
			Subsystem02: reactor.RequiredStatus_UpOrCeded,
			Subsystem03: reactor.RequiredStatus_UpOrCeded,
		}
		r.Func = func(r *reactor.Reactor, m *reactor.Msg) {
			switch m.Signal {
			case reactor.Signal_End:
				r.Tracer.LogEvent("22").SendMsg("end")

			case reactor.Signal_IsUp, reactor.Signal_Cede:
				r.Tracer.LogEvent("24").SendMsg("ready signal")
				status := r.GetStatus_InsideReactorFuncOnly()
				r.Tracer.LogEvent("29").
					Strs("fresh", status.Fresh).
					Strs("ready", status.Up).
					Strs("ceded", status.Ceded).
					SendMsg(m.Name + " signalled")

				// this is common idiom
				if status.CanReact {
					r.Tracer.LogEvent("45").SendMsg("READINESS ACHIEVED - DO THE NEXT STEP")
					r.Done = true
					// do something
				} else if status.NoFreshLeft {
					r.Tracer.LogEvent("47").SendMsg("NO FRESH LEFT - DO THE NEXT STEP")
					r.Done = true
					// do something
				}
			}
		}
		r.StartFresh()

		fmt.Println("STARTED")

		time.Sleep(time.Second * 1)
		r.Send(&reactor.Msg{Signal: reactor.Signal_End})

		time.Sleep(time.Second * 1)
		r.End()
		r.End() // this must be OK

		// panic test
		r.SendStatus_Up(nil, Subsystem01)

		return
	}).Catch(func(e *controlflow.Event) {
		fmt.Println(e)
		require.Equal(t, `{: FORWARDTRACE {: ETC [1 1 2], TAGS [1e0194 R1]}, BACKTRACE { [eae98cf2227f "send on closed channel"] [01967ad5b672bff68807 "reactor is not inited or has been ended"]}}`, e.String())
	}).OK(func() {
		t.FailNow()
	})
}
