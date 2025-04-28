package delayed

import (
	"context"
	"time"

	"github.com/rusriver/nutz/logger/loggerapi"
)

type Delayed struct {
	Name       string
	Context    context.Context
	Logger     loggerapi.ILogger
	Duration   time.Duration
	F          func(d *Delayed)
	F_canceled func(d *Delayed)
}

func Start(d *Delayed) {
	go func() {
		select {
		case <-d.Context.Done():
			if d.F_canceled != nil {
				d.F_canceled(d)
			}
		case <-time.After(d.Duration):
			if d.F != nil {
				d.F(d)
			}
		}
	}()
}
