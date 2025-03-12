package logger

import (
	"io"

	"github.com/rusriver/nutz/logger/loggerapi"
)

type Interceptor struct {
	rctx *loggerapi.RootContext
	e    loggerapi.IEvent
}

var _ io.Writer = (*Interceptor)(nil)

func (i *Interceptor) Write(p []byte) (int, error) {
	if i.rctx.EventOnWireHook != nil {
		p2 := append(p[:0:0], p...)
		// p2 := make([]byte, len(p))
		// copy(p2, p)
		i.rctx.EventOnWireHook(i.e, p2)
	}
	return i.rctx.MainWriter.Write(p)
}
