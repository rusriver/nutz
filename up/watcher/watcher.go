package watcher

import (
	"context"
	"time"

	loggeterrapi "github.com/rusriver/nutz/loggeterr-api"
)

type Watcher struct {
	Name         string
	Tracer       loggeterrapi.ITracer
	Context      context.Context
	PeriodicFunc func(w *Watcher) // optional
	Period       time.Duration    // optional (can be zero)
	Timeout      time.Duration
	TimeoutFunc  func(w *Watcher) // required
	chUp         chan *cmd
}

type cmd struct {
	Command string
}

func (w *Watcher) Start() {
	if w.chUp != nil {
		// already started
		return
	}
	w.Tracer.LogEvent("e325ac80ec09").SendMsgf("the watcher %v starting", w.Name)
	w.chUp = make(chan *cmd, 0xF)
	go func() {
		var ticker *time.Ticker
		if w.Period > 0 {
			ticker = time.NewTicker(w.Period)
		} else {
			ticker = &time.Ticker{
				C: make(<-chan time.Time),
			}
		}
		after_timeout := time.After(w.Timeout)
		for {
			select {
			case <-w.Context.Done():
				return
			case c := <-w.chUp:
				switch c.Command {
				case "stop_done":
					w.Tracer.LogEvent("5ad2754f035d").SendMsgf("the watcher %v reached STOP/DONE state (OK)", w.Name)
					return
				}
			case <-ticker.C:
				w.Tracer.LogEvent("c401d1be793c").SendMsgf("the watcher %v still waiting...", w.Name)
				if w.PeriodicFunc != nil {
					w.PeriodicFunc(w)
				}
			case <-after_timeout:
				w.Tracer.LogEvent("b7d60e67e04c").SendMsgf("the watcher %v timed out", w.Name)
				w.TimeoutFunc(w)
			}
		}
	}()
}

func (w *Watcher) Stop_Done() {
	if w.chUp == nil {
		// already stopped
		return
	}
	w.chUp <- &cmd{Command: "stop_done"}
}
