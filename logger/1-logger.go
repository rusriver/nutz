package logger

import (
	"errors"
	"io"
	"os"
	"sync"

	"github.com/lithammer/shortuuid/v4"
	"github.com/rs/zerolog"
	"github.com/rusriver/configdb/v3"
	"github.com/rusriver/nutz/logger/loggerapi"
)

// Single-thread only! Use sub-loggers if want to use it concurrently
type Logger struct {
	RootContext  *loggerapi.RootContext
	IsInactive   bool // per logger
	UserObject   any
	TraceContext *loggerapi.TraceContext

	ZeroLogger *zerolog.Logger

	mu          sync.Mutex // only some functions use it
	accumulator *Accumulator
	interceptor *Interceptor
}

var _ loggerapi.ILogger = (*Logger)(nil)

func New(xuid string, zl *zerolog.Logger, optFuncs ...func(rctx *loggerapi.RootContext)) (logger *Logger) {
	rctx := &loggerapi.RootContext{
		MsgtagKey: "msgtags",
		XUIdKey:   "xuid",
		ETCKey:    "etc",
	}
	for _, f := range optFuncs {
		f(rctx)
	}
	if rctx.MainWriter == nil {
		rctx.MainWriter = os.Stderr
	}
	rctx.XUId = xuid
	if len(rctx.XUId) == 0 {
		rctx.XUId = shortuuid.New()[:8]
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
	interceptor := &Interceptor{rctx: rctx}
	if zl == nil {
		zlv := zerolog.New(interceptor).With().Timestamp().Logger()
		zl = &zlv
	}
	logger = &Logger{
		RootContext: rctx,
		ZeroLogger:  zl,
		TraceContext: &loggerapi.TraceContext{
			ETC: loggerapi.EventTracingCounters{0},
		},
		interceptor: interceptor,
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

func (l *Logger) SubLoggerInitChain(msgtag string) loggerapi.ISubLoggerInitChain {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.TraceContext.ETC.Inc()

	zerologContext := l.ZeroLogger.With()
	slich := &subLoggerInitChain{
		zerologContext: &zerologContext,
		traceContext: &loggerapi.TraceContext{
			ETC:     l.TraceContext.ETC.BranchSub(),
			Msgtags: l.TraceContext.Msgtags.With(msgtag),
		},
		rctx:        l.RootContext,
		uo:          l.UserObject,
		interceptor: l.interceptor,
	}
	return slich
}

func appendTracingData(e *Event, msgtag string) {
	e.ParentLogger.mu.Lock()
	defer e.ParentLogger.mu.Unlock()

	k := e.ParentLogger.RootContext.XUIdKey
	if len(k) != 0 {
		e.ZerologEvent = e.ZerologEvent.Str(k, e.ParentLogger.RootContext.XUId)
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

func (l *Logger) SwitchOnTheOfflineAccumulation() loggerapi.ILogger {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.accumulator = NewAccumulator()

	rctx := l.RootContext.Clone()

	rctx.EventOnWireHook = func(e loggerapi.IEventInHook, p []byte) {
		l.accumulator.chUp <- p
	}
	rctx.MainWriter = io.Discard

	l.interceptor = &Interceptor{rctx: rctx}

	zlv := l.ZeroLogger.Output(l.interceptor)
	l.ZeroLogger = &zlv

	return l
}

func (l *Logger) FinishAccumulationAndFetch() (llines []*configdb.Root) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.accumulator == nil {
		panic("401d5fc2a82e553f0eab")
	}
	return l.accumulator.FinishAccumulationAndFetch()
}

func (l *Logger) RestartTheAccumulation() {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.accumulator = NewAccumulator()

	l.RootContext.EventOnWireHook = func(e loggerapi.IEventInHook, p []byte) {
		l.accumulator.chUp <- p
	}
}

func (l *Logger) GetUserObject() (uo any) {
	l.mu.Lock()
	defer l.mu.Unlock()

	return l.UserObject
}

func (l *Logger) GetTraceContext() (tc *loggerapi.TraceContext) {
	return l.TraceContext
}
