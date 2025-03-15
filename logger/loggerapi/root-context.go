package loggerapi

import "io"

type RootContext struct {
	MainWriter       io.Writer
	MsgtagKey        string
	XUIdKey          string // XU logger root id
	ETCKey           string // event tracing counters
	IfSendHook       func(e IEventInHook) (doSend bool)
	EventOnWireHook  func(e IEventInHook, p []byte)
	ActivationHook   func(e IEventInHook) (inactivate bool) // If the ActicationHook is set, it is called, and it can reactivate the event.
	NewSubLoggerHook func(tc *TraceContext)
	SpanBeginHook    func(tc *TraceContext)
	SpanEndHook      func(tc *TraceContext)
	XUId             string
}

type IEventInHook interface {
	GetTraceContext() *TraceContext
	GetUserObject() (uo any)
}

func (rctx0 *RootContext) Clone() *RootContext {
	rctx2 := &RootContext{
		MainWriter:       rctx0.MainWriter,
		MsgtagKey:        rctx0.MsgtagKey,
		XUIdKey:          rctx0.XUIdKey,
		ETCKey:           rctx0.ETCKey,
		IfSendHook:       rctx0.IfSendHook,
		EventOnWireHook:  rctx0.EventOnWireHook,
		ActivationHook:   rctx0.ActivationHook,
		NewSubLoggerHook: rctx0.NewSubLoggerHook,
		SpanBeginHook:    rctx0.SpanBeginHook,
		SpanEndHook:      rctx0.SpanEndHook,
		XUId:             rctx0.XUId,
	}
	return rctx2
}
