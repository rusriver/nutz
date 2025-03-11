package logger

import (
	"errors"
	"os"
	"sync"

	"github.com/lithammer/shortuuid/v4"
	"github.com/rs/zerolog"
	"github.com/rusriver/nutz/logger/loggerapi"
)

// Single-thread only! Use sub-loggers if want to use it concurrently
type Logger struct {
	RootContext  *loggerapi.RootContext
	IsInactive   bool // per logger
	UserObject   any
	TraceContext *loggerapi.TraceContext

	ZeroLogger *zerolog.Logger
	mu         sync.Mutex // only some functions use it
	_XUId      string
}

var _ loggerapi.ILogger = (*Logger)(nil)

func New(xuid string, zl *zerolog.Logger, optFuncs ...func(rctx *loggerapi.RootContext)) (logger *Logger) {
	rctx := &loggerapi.RootContext{
		PanicOnMisuse: true,
		MsgtagKey:     "msgtags",
		XUIdKey:       "xuid",
		ETCKey:        "etc",
	}
	for _, f := range optFuncs {
		f(rctx)
	}
	if rctx.MainWriter == nil {
		rctx.MainWriter = os.Stderr
	}
	{
		zerolog.TimestampFieldName = "ts"
		zerolog.TimeFieldFormat = "20060102-150405"
		zerolog.LevelTraceValue = "TRACE"
		zerolog.LevelDebugValue = "DIAGN"
		zerolog.LevelInfoValue = "INFO"
		zerolog.LevelWarnValue = "WARN"
		zerolog.LevelErrorValue = "ERROR"
		zerolog.LevelFatalValue = "FATAL"
		zerolog.LevelPanicValue = "PANIC"
	}
	if zl == nil {
		i := &Interceptor{rctx: rctx}
		zlv := zerolog.New(i).With().Timestamp().Logger()
		zl = &zlv
	}
	logger = &Logger{
		RootContext: rctx,
		ZeroLogger:  zl,
		TraceContext: &loggerapi.TraceContext{
			ETC: loggerapi.EventTracingCounters{0},
		},
		_XUId: xuid,
	}
	if len(logger._XUId) == 0 {
		logger._XUId = shortuuid.New()[:8]
	}
	return logger
}

func (l *Logger) TraceEvent(msgtag string) loggerapi.IEvent {
	e := &Event{
		IsInactive:   l.IsInactive,
		TheLevel:     zerolog.TraceLevel.String(),
		ParentLogger: l,
		ZerologEvent: l.ZeroLogger.Trace(),
	}
	appendTracingData(e, msgtag)
	return e
}

func (l *Logger) DiagnosticEvent(msgtag string) loggerapi.IEvent {
	e := &Event{
		IsInactive:   l.IsInactive,
		TheLevel:     zerolog.DebugLevel.String(),
		ParentLogger: l,
		ZerologEvent: l.ZeroLogger.Debug(),
	}
	appendTracingData(e, msgtag)
	return e
}

func (l *Logger) InfoEvent(msgtag string) loggerapi.IEvent {
	e := &Event{
		IsInactive:   l.IsInactive,
		TheLevel:     zerolog.InfoLevel.String(),
		ParentLogger: l,
		ZerologEvent: l.ZeroLogger.Info(),
	}
	appendTracingData(e, msgtag)
	return e
}

func (l *Logger) WarnEvent(msgtag string) loggerapi.IEvent {
	e := &Event{
		IsInactive:        l.IsInactive,
		TheLevel:          zerolog.WarnLevel.String(),
		ParentLogger:      l,
		ZerologEvent:      l.ZeroLogger.Warn(),
		TheTitle:          "WARNING", // the default title
		theTitleIsDefault: true,
	}
	appendTracingData(e, msgtag)
	return e
}

func (l *Logger) ErrorEvent(msgtag string) loggerapi.IEvent {
	e := &Event{
		IsInactive:        l.IsInactive,
		TheLevel:          zerolog.ErrorLevel.String(),
		ParentLogger:      l,
		ZerologEvent:      l.ZeroLogger.Error(),
		TheTitle:          "ERROR", // the default title
		theTitleIsDefault: true,
	}
	appendTracingData(e, msgtag)
	return e
}

func (l *Logger) ErrEvent(err error, msgtag string) loggerapi.IEvent {
	if err == nil {
		err = errors.New("<nil>")
	}
	e := &Event{
		IsInactive:        l.IsInactive,
		TheLevel:          zerolog.ErrorLevel.String(),
		ParentLogger:      l,
		ZerologEvent:      l.ZeroLogger.Err(err),
		TheTitle:          "ERROR", // the default title
		theTitleIsDefault: true,
	}
	appendTracingData(e, msgtag)
	return e
}

func (l *Logger) SubLoggerInitChain(msgtag string) loggerapi.IEvent {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.TraceContext.ETC.Inc()

	zerologContext := l.ZeroLogger.With()
	e := &Event{
		subLoggerInitChain: &subLoggerInitChain{
			zerologContext: &zerologContext,
			msgtags:        l.TraceContext.Msgtags.With(msgtag),
			_ETC:           l.TraceContext.ETC.BranchSub(),
			newRctx:        l.RootContext,
			newUserObject:  l.UserObject,
		},
		ParentLogger: l,
	}
	return e
}

func (l *Logger) SetLevel(level zerolog.Level) loggerapi.ILogger {
	return l
}

func appendTracingData(e *Event, msgtag string) {
	e.ParentLogger.mu.Lock()
	defer e.ParentLogger.mu.Unlock()

	k := e.ParentLogger.RootContext.XUIdKey
	if len(k) != 0 {
		e.ZerologEvent = e.ZerologEvent.Str(k, e.ParentLogger._XUId)
	}

	k = e.ParentLogger.RootContext.ETCKey
	if len(k) != 0 {
		e.ParentLogger.TraceContext.ETC.Inc()
		e.ZerologEvent = e.ZerologEvent.Ints(k, e.ParentLogger.TraceContext.ETC)
	}

	e.TheMsgtags = e.ParentLogger.TraceContext.Msgtags.With(msgtag)

	// call the (re-)activation hook, if any
	hook := e.ParentLogger.RootContext.ActivationHook
	if hook != nil {
		e.IsInactive = hook(e)
	}
	k = e.ParentLogger.RootContext.MsgtagKey
	if len(k) != 0 {
		e.ZerologEvent.Strs(k, e.TheMsgtags)
	}
}
