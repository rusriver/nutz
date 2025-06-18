package loggeterrapi

import (
	"io"

	"github.com/rusriver/nutz/loggeterr-api/trace"
)

type RootContext struct {
	MainWriter          io.Writer
	ForwardtraceTagsKey string
	ForwardtraceETCKey  string // event tracing counters
	IfSendHook          func(e IEventInHook) (doSend bool)
	EventOnWireHook     func(e IEventInHook, p []byte)
	ActivationHook      func(e IEventInHook) (inactivate bool) // If the ActicationHook is set, it is called, and it can reactivate the event.
	NewSubLoggerHook    func(ft *trace.Forwardtrace)
	SpanBeginHook       func(ft *trace.Forwardtrace)
	SpanEndHook         func(ft *trace.Forwardtrace)
}

type IEventInHook interface {
	GetForwardtrace() *trace.Forwardtrace
	GetUserObject() (uo any)
}

func (rctx0 *RootContext) Clone() *RootContext {
	rctx2 := &RootContext{
		MainWriter:          rctx0.MainWriter,
		ForwardtraceTagsKey: rctx0.ForwardtraceTagsKey,
		ForwardtraceETCKey:  rctx0.ForwardtraceETCKey,
		IfSendHook:          rctx0.IfSendHook,
		EventOnWireHook:     rctx0.EventOnWireHook,
		ActivationHook:      rctx0.ActivationHook,
		NewSubLoggerHook:    rctx0.NewSubLoggerHook,
		SpanBeginHook:       rctx0.SpanBeginHook,
		SpanEndHook:         rctx0.SpanEndHook,
	}
	return rctx2
}
