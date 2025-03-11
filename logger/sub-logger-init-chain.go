package logger

import (
	"github.com/rs/zerolog"
	"github.com/rusriver/nutz/logger/loggerapi"
)

type subLoggerInitChain struct {
	zerologContext *zerolog.Context
	isInactive     bool // for new sub-logger
	msgtags        loggerapi.Msgtags
	_ETC           loggerapi.EventTracingCounters
	newRctx        *loggerapi.RootContext
	newUserObject  any
}
