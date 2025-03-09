package logger

import (
	"github.com/rs/zerolog"
	"github.com/rusriver/nutz/logger/loggerapi"
)

type subLoggerInitChain struct {
	zerologContext     *zerolog.Context
	perLoggerSubsystem string
	isInactive         bool // for new sub-logger
	msgtags            *loggerapi.Msgtags
}

type RootContext struct {
	PanicOnMisuse  bool
	MsgtagKey      string
	LRIdKey        string // logger root id
	ETCKey         string // event tracing counters
	OnSendHook     func(e *Event) (doSend bool)
	ActivationHook func(e *Event) (inactivate bool) // If the ActicationHook is set, it is called, and it can reactivate the event.
	RootId         string
}
