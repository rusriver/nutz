package reactor

import (
	"context"
	"sync"

	"github.com/rusriver/nutz/controlflow"
	loggeterrapi "github.com/rusriver/nutz/loggeterr-api"
	"github.com/rusriver/nutz/loggeterr-api/trace"
	upcommon "github.com/rusriver/nutz/up/up-common"
)

type Reactor struct {
	Name          string // the name is required, since it's used in a user-side tracer when calling RecursiveWindUp()
	Context       context.Context
	Tracer        loggeterrapi.ITracer
	NeedsPattern  SubsystemsPattern // the required dependencies pattern, map[<dep_name>]<required_status>
	Func          func(r *Reactor, msg *Msg)
	Func_CanReact func(r *Reactor, msg *Msg, status *Status)
	subsystems    Subsystems // must be unexported, because of SRP
	Done          bool
	chUp          chan *Msg
	opts          *Options
}

type Options struct {
	ChUpSize int
}

func (r *Reactor) StartFresh(of ...func(opts *Options)) *Reactor {
	if len(r.Name) == 0 {
		r.Tracer.RaiseCause(&trace.BacktracePoint{"794b7eafe45c", "reactor name is empty"})
	}

	r.subsystems = make(Subsystems)
	for name := range r.NeedsPattern {
		r.subsystems[name] = &Subsystem{
			Name:   name,
			Status: Status_Fresh,
		}
	}

	r.opts = &Options{
		ChUpSize: 1,
	}
	for _, f := range of {
		f(r.opts)
	}

	r.chUp = make(chan *Msg, r.opts.ChUpSize)
	go func() {
		for {
			select {
			case <-r.Context.Done():
				return
			case m := <-r.chUp:
				switch m.Signal {
				case Signal_CallTheInjectFunc:
					m.InjectFunc(r)
				case Signal_IsUp, Signal_Cede, Signal_IsFailed:
					sub := r.subsystems[m.Name]
					if sub == nil {
						sub = &Subsystem{Name: m.Name}
					}
					switch m.Signal {
					case Signal_IsUp:
						sub.Status = Status_Up
					case Signal_IsFailed:
						sub.Status = Status_Failed
					case Signal_Cede:
						sub.Status = Status_Ceded
					}
					sub.Itself = m.OTDG
					r.subsystems[m.Name] = sub

					if r.Func_CanReact != nil {
						status := r.GetStatus_InsideReactorFuncOnly()
						if status.CanReact {
							r.Func_CanReact(r, m, status)
						}
					}
				}

				if r.Func != nil {
					r.Func(r, m)
				}

				switch m.Signal {
				case Signal_End:
					close(r.chUp)
					return
				}
			}
		}
	}()

	return r
}

func (r *Reactor) Send(m *Msg) {
	defer func() {
		if x := recover(); x != nil {
			event := controlflow.NewEvent(r.Tracer.GetForwardtrace(), &trace.BacktracePoint{"eae98cf2227f", x})
			event.Reraise(trace.BacktracePointf("01967ad5b672bff68807", "reactor is not inited or has been ended"))
		}
	}()
	r.chUp <- m
}

func (r *Reactor) SendStatus_Up(otdg upcommon.IRecursiveWindUp, name string) {
	r.Send(&Msg{Signal_IsUp, name, otdg, nil})
}

func (r *Reactor) SendStatus_Fail(otdg upcommon.IRecursiveWindUp, name string) {
	r.Send(&Msg{Signal_IsFailed, name, otdg, nil})
}

func (r *Reactor) SendStatus_Cede(otdg upcommon.IRecursiveWindUp, name string) {
	r.Send(&Msg{Signal_Cede, name, otdg, nil})
}

// Prefer not to use this when reactor achieves readiness - what if some will send up signals after
// the reactor already achieved readiness and closed? It will panic, and that's not what you may want.
// Use Done flag to make sure it reacts only once.
func (r *Reactor) End() {
	defer func() {
		recover()
	}()
	r.chUp <- &Msg{Signal: Signal_End}
}

type Status struct {
	Fresh             []string
	Up                []string
	Ceded             []string
	NoFreshLeft       bool
	ReadinessAchieved bool
	CanReact          bool // it is set if ReadinessAchieved && !Done
}

// Please: use only from within the reactor func
func (r *Reactor) GetStatus_InsideReactorFuncOnly() (s *Status) {
	s = &Status{
		NoFreshLeft: true, // the default
	}
	for name, sub := range r.subsystems {
		switch sub.Status {
		case Status_Fresh:
			s.Fresh = append(s.Fresh, name)
			s.NoFreshLeft = false
		case Status_Up:
			s.Up = append(s.Up, name)
		case Status_Ceded:
			s.Ceded = append(s.Ceded, name)
		}
	}
	checks := 0
	for name, required_status := range r.NeedsPattern {
		switch required_status {
		case RequiredStatus_Up:
			if sub, ok := r.subsystems[name]; ok {
				if sub.Status == Status_Up {
					checks++
				}
			}
		case RequiredStatus_UpOrCeded:
			if sub, ok := r.subsystems[name]; ok {
				if sub.Status == Status_Up || sub.Status == Status_Ceded {
					checks++
				}
			}
		}
	}
	if checks == len(r.NeedsPattern) {
		s.ReadinessAchieved = true
	}
	if s.ReadinessAchieved && !r.Done {
		s.CanReact = true
	}
	return
}

func (r *Reactor) CallTheInjectFunc(f func(r *Reactor)) {
	r.Send(&Msg{Signal: Signal_CallTheInjectFunc, InjectFunc: f})
}

// Be careful if your reactor dependency graph has diamond topology (is not a DAG);
// In that case, better implement explicit forward action plan to wind up your system.
func (r *Reactor) RecursiveWindUp(tracer loggeterrapi.ITracer) {
	tracer = tracer.NSTIC(r.Name).ITracer()
	tracer.InfoEvent("184").SendMsg("rec-windup-begin")
	defer tracer.InfoEvent("185").SendMsg("rec-windup-end")
	var subs []*Subsystem
	wg := sync.WaitGroup{}
	wg.Add(1)
	r.CallTheInjectFunc(func(r *Reactor) {
		for _, sub := range r.subsystems {
			switch sub.Status {
			case Status_Up, Status_Ceded, Status_Failed:
				subs = append(subs, sub)
			}
		}
		wg.Done()
	})
	wg.Wait()
	for _, sub := range subs {
		if sub.Itself != nil {
			wg.Add(1)
			go func() {
				sub.Itself.RecursiveWindUp(tracer)
				wg.Done()
			}()
		}
	}
	wg.Wait()
}
