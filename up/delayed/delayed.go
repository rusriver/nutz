package delayed

import (
	"context"
	"time"

	loggeterrapi "github.com/rusriver/nutz/loggeterr-api"
)

type Delayed struct {
	Context       context.Context
	Tracer        loggeterrapi.ITracer
	Duration      time.Duration
	Func          func(d *Delayed)
	Func_canceled func(d *Delayed)
}

func Start(d *Delayed) {
	go func() {
		select {
		case <-d.Context.Done():
			if d.Func_canceled != nil {
				d.Func_canceled(d)
			}
		case <-time.After(d.Duration):
			if d.Func != nil {
				d.Func(d)
			}
		}
	}()
}
