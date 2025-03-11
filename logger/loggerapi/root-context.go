package loggerapi

import "io"

type RootContext struct {
	MainWriter      io.Writer
	PanicOnMisuse   bool
	MsgtagKey       string
	XUIdKey         string // XU logger root id
	ETCKey          string // event tracing counters
	IfSendHook      func(e IEvent) (doSend bool)
	EventOnWireHook func(p []byte)
	ActivationHook  func(e IEvent) (inactivate bool) // If the ActicationHook is set, it is called, and it can reactivate the event.
}
