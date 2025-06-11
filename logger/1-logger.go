package logger

import (
	"errors"
	"io"
	"os"
	"sync"

	"github.com/rs/zerolog"
	"github.com/rusriver/configdb/v3"
	"github.com/rusriver/nutz/logger/loggerapi"
	"github.com/rusriver/nutz/logger/loggerapi/accumapi"
)

// Single-thread only! Use sub-loggers if want to use it concurrently
type Logger struct {
	RootContext  *loggerapi.RootContext
	IsInactive   bool
	UserObject   any
	TraceContext *loggerapi.TraceContext

	ZeroLogger *zerolog.Logger

	mu          sync.Mutex // only some functions use it
	accumulator *Accumulator
	interceptor *Interceptor
}

var _ loggerapi.ILogger = (*Logger)(nil)
var _ accumapi.ILogsAccumulator = (*Logger)(nil)

func New(zl *zerolog.Logger, optFuncs ...func(rctx *loggerapi.RootContext)) (logger loggerapi.ILogger) {
	rctx := &loggerapi.RootContext{
		MsgtagKey: "msgtags",
		ETCKey:    "etc",
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

func appendTracingData(e *Event, tag string) {
	e.ParentLogger.mu.Lock()
	defer e.ParentLogger.mu.Unlock()

	tc := e.ParentLogger.TraceContext
	rc := e.ParentLogger.RootContext

	k := rc.ETCKey
	if len(k) != 0 {
		tc.ETC.Inc()
		e.ZerologEvent = e.ZerologEvent.Ints(k, tc.ETC)
	}

	tc.Msgtags = append(tc.Msgtags, tag)
	// The "pop" operation is at the e.Send(), and it's race-unsafe intentionally -
	// because if we e.g. lock the mutex, we have the risk to dead-lock incorrect user
	// program. The choise is made to prefer incorrect tagging, but not dead-locking,
	// in that case.

	// call the (re-)activation hook, if any
	hook := rc.ActivationHook
	if hook != nil {
		e.IsInactive = hook(e)
		tc.IsInactive = e.IsInactive
	}
	k = rc.MsgtagKey
	if len(k) != 0 {
		e.ZerologEvent.Strs(k, tc.Msgtags)
	}
}

func (l *Logger) NewSubLoggerInitChain(msgtag string) loggerapi.ISubLoggerInitChain {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.TraceContext.ETC.Inc()

	zerologContext := l.ZeroLogger.With()
	slich := &subLoggerInitChain{
		zerologContext: &zerologContext,
		traceContext: &loggerapi.TraceContext{
			ETC:     l.TraceContext.ETC.CloneBranchSub(),
			Msgtags: l.TraceContext.Msgtags.CloneWith(msgtag),
		},
		rctx:        l.RootContext,
		uo:          l.UserObject,
		interceptor: l.interceptor,
	}

	h := l.RootContext.NewSubLoggerHook
	if h != nil {
		h(slich.traceContext)
	}

	return slich
}

func (l *Logger) NSLIC(msgtag string) loggerapi.ISubLoggerInitChain {
	return l.NewSubLoggerInitChain(msgtag)
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

func (l *Logger) SpanBegin(tag string) loggerapi.ILogger {
	l.mu.Lock()
	defer l.mu.Unlock()

	tc := l.TraceContext

	tc.ETC = append(tc.ETC, 0)
	tc.Msgtags = append(tc.Msgtags, tag)

	h := l.RootContext.SpanBeginHook
	if h != nil {
		h(tc)
	}

	return l
}

func (l *Logger) SpanEnd() {
	l.mu.Lock()
	defer l.mu.Unlock()

	tc := l.TraceContext

	h := l.RootContext.SpanEndHook
	if h != nil {
		h(tc)
	}

	tc.ETC = tc.ETC[:len(tc.ETC)-1]
	tc.Msgtags = tc.Msgtags[:len(tc.Msgtags)-1]
}
