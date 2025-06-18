package reactor

import (
	"context"

	"github.com/rusriver/nutz/controlflow"
	"github.com/rusriver/nutz/logger/loggerapi"
)

type Reactor struct {
	Name       string
	Context    context.Context
	Logger     loggerapi.ILogger
	F          func(r *Reactor, msg *Msg)
	subsystems subsystems
	chUp       chan *Msg
	opts       *Options
}

type Options struct {
	ChUpSize int
}

func Start(r *Reactor, _subsystems []string, of ...func(opts *Options)) *Reactor {
	r.subsystems = make(subsystems)
	for _, ss := range _subsystems {
		r.subsystems[ss] = SubsystemStatus_Fresh
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
				case Signal_Ready:
					r.subsystems[m.PId] = SubsystemStatus_Ready // allows to report additional subs as well
				case Signal_Cede:
					r.subsystems[m.PId] = SubsystemStatus_Ceded
				}

				if r.F == nil {
					controlflow.Raise1("01967ae6cdb08172b74c reactor func is nil")
				}
				r.F(r, m)

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
			controlflow.Reraise(x, "01967ad5b672bff68807 reactor is not inited or has been ended")
		}
	}()
	r.chUp <- m
}

func (r *Reactor) End() {
	defer func() {
		recover()
	}()
	r.chUp <- &Msg{Signal: Signal_End}
}

type Status struct {
	Fresh   []string
	Ready   []string
	Ceded   []string
	AllDone bool
}

// use only from within the reactor func
func (r *Reactor) Status() (s *Status) {
	s = &Status{
		AllDone: true,
	}
	for ss, status := range r.subsystems {
		switch status {
		case SubsystemStatus_Fresh:
			s.Fresh = append(s.Fresh, ss)
			s.AllDone = false
		case SubsystemStatus_Ready:
			s.Ready = append(s.Ready, ss)
		case SubsystemStatus_Ceded:
			s.Ceded = append(s.Ceded, ss)
		}
	}
	return
}
