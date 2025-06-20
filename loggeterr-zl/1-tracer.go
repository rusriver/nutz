package loggeterrzl

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/rs/zerolog"
	// "github.com/rusriver/configdb/v3"
	"github.com/rusriver/configdb/v3"
	"github.com/rusriver/nutz/controlflow"
	loggeterrapi "github.com/rusriver/nutz/loggeterr-api"
	"github.com/rusriver/nutz/loggeterr-api/accumapi"
	"github.com/rusriver/nutz/loggeterr-api/trace"
)

// Single-thread only! Use sub-loggers if want to use it concurrently
type Tracer struct {
	RootContext  *loggeterrapi.RootContext
	IsInactive   bool
	UserObject   any
	Forwardtrace *trace.Forwardtrace

	ZeroLogger *zerolog.Logger

	mu          sync.Mutex // only some functions use it
	accumulator *Accumulator
	interceptor *Interceptor
}

var _ loggeterrapi.ITracer = (*Tracer)(nil)
var _ accumapi.ILogsAccumulator = (*Tracer)(nil)

func New(zl *zerolog.Logger, optFuncs ...func(rctx *loggeterrapi.RootContext)) (tracer loggeterrapi.ITracer) {
	rctx := &loggeterrapi.RootContext{
		ForwardtraceTagsKey: "tags",
		ForwardtraceETCKey:  "etc",
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
	tracer = &Tracer{
		RootContext: rctx,
		ZeroLogger:  zl,
		Forwardtrace: &trace.Forwardtrace{
			ETC: trace.EventTracingCounters{0},
		},
		interceptor: interceptor,
	}
	return tracer
}

func (t *Tracer) TraceEvent(msgtag string) loggeterrapi.ILogEvent {
	e := &LogEvent{
		IsInactive:   t.IsInactive,
		TheLevel:     zerolog.TraceLevel.String(),
		ParentLogger: t,
		ZerologEvent: t.ZeroLogger.Trace(),
	}
	appendTracingData(e, msgtag)
	return e
}

func (t *Tracer) DiagnosticEvent(msgtag string) loggeterrapi.ILogEvent {
	e := &LogEvent{
		IsInactive:   t.IsInactive,
		TheLevel:     zerolog.DebugLevel.String(),
		ParentLogger: t,
		ZerologEvent: t.ZeroLogger.Debug(),
	}
	appendTracingData(e, msgtag)
	return e
}

func (t *Tracer) InfoEvent(msgtag string) loggeterrapi.ILogEvent {
	e := &LogEvent{
		IsInactive:   t.IsInactive,
		TheLevel:     zerolog.InfoLevel.String(),
		ParentLogger: t,
		ZerologEvent: t.ZeroLogger.Info(),
	}
	appendTracingData(e, msgtag)
	return e
}

func (t *Tracer) WarnEvent(msgtag string) loggeterrapi.ILogEvent {
	e := &LogEvent{
		IsInactive:        t.IsInactive,
		TheLevel:          zerolog.WarnLevel.String(),
		ParentLogger:      t,
		ZerologEvent:      t.ZeroLogger.Warn(),
		TheTitle:          "WARNING", // the default title
		theTitleIsDefault: true,
	}
	appendTracingData(e, msgtag)
	return e
}

func (t *Tracer) ErrorEvent(msgtag string) loggeterrapi.ILogEvent {
	e := &LogEvent{
		IsInactive:        t.IsInactive,
		TheLevel:          zerolog.ErrorLevel.String(),
		ParentLogger:      t,
		ZerologEvent:      t.ZeroLogger.Error(),
		TheTitle:          "ERROR", // the default title
		theTitleIsDefault: true,
	}
	appendTracingData(e, msgtag)
	return e
}

func (t *Tracer) SmartErrEvent(err any, msgtag string) loggeterrapi.ILogEvent {
	if err == nil {
		err = errors.New("<nil>")
	}
	e := &LogEvent{
		IsInactive:        t.IsInactive,
		TheLevel:          zerolog.ErrorLevel.String(),
		ParentLogger:      t,
		TheTitle:          "ERROR", // the default title
		theTitleIsDefault: true,
	}
	switch ev := err.(type) {
	case *controlflow.Event:
		e.ZerologEvent = t.ZeroLogger.Error()
		bb, _ := json.Marshal(ev)
		e.ZerologEvent = t.ZeroLogger.Error().RawJSON("error-controlflow-event", bb)
	case *trace.BacktracePoint:
		bb, _ := json.Marshal(ev)
		e.ZerologEvent = t.ZeroLogger.Error().RawJSON("error-backtracepoint", bb)
	case error:
		e.ZerologEvent = t.ZeroLogger.Err(ev)
	default:
		e.ZerologEvent = t.ZeroLogger.Error().Str("error-unknown", fmt.Sprintf("%+v", ev))
	}
	appendTracingData(e, msgtag)
	return e
}

func (t *Tracer) LogEvent(msgtag string) loggeterrapi.ILogEvent {
	e := &LogEvent{
		IsInactive:   t.IsInactive,
		TheLevel:     zerolog.InfoLevel.String(),
		ParentLogger: t,
		ZerologEvent: t.ZeroLogger.Info(),
	}
	appendTracingData(e, msgtag)
	return e
}

func appendTracingData(e *LogEvent, tag string) {
	e.ParentLogger.mu.Lock()
	defer e.ParentLogger.mu.Unlock()

	tc := e.ParentLogger.Forwardtrace
	rc := e.ParentLogger.RootContext

	k := rc.ForwardtraceETCKey
	if len(k) != 0 {
		tc.ETC.Inc()
		e.ZerologEvent = e.ZerologEvent.Ints(k, tc.ETC)
	}

	tc.Tags = append(tc.Tags, tag)
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
	k = rc.ForwardtraceTagsKey
	if len(k) != 0 {
		e.ZerologEvent.Strs(k, tc.Tags)
	}
}

func (t *Tracer) NewSubTracerInitChain(msgtag string) loggeterrapi.ISubTracerInitChain {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.Forwardtrace.ETC.Inc()

	zerologContext := t.ZeroLogger.With()
	slich := &subTracerInitChain{
		zerologContext: &zerologContext,
		forwardtrace: &trace.Forwardtrace{
			ETC:  t.Forwardtrace.ETC.CloneBranchSub(),
			Tags: t.Forwardtrace.Tags.CloneWith(msgtag),
		},
		rctx:        t.RootContext,
		uo:          t.UserObject,
		interceptor: t.interceptor,
	}

	h := t.RootContext.NewSubLoggerHook
	if h != nil {
		h(slich.forwardtrace)
	}

	return slich
}

func (t *Tracer) NSTIC(msgtag string) loggeterrapi.ISubTracerInitChain {
	return t.NewSubTracerInitChain(msgtag)
}

func (t *Tracer) SwitchOnTheOfflineAccumulation() loggeterrapi.ITracer {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.accumulator = NewAccumulator(t)

	rctx := t.RootContext.Clone()

	rctx.EventOnWireHook = func(e loggeterrapi.IEventInHook, p []byte) {
		t.accumulator.chUp <- p
	}
	rctx.MainWriter = io.Discard

	t.interceptor = &Interceptor{rctx: rctx}

	zlv := t.ZeroLogger.Output(t.interceptor)
	t.ZeroLogger = &zlv

	return t
}

func (t *Tracer) FinishAccumulationAndFetch() (llines []*configdb.Root) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.accumulator == nil {
		panic("401d5fc2a82e553f0eab")
	}
	return t.accumulator.FinishAccumulationAndFetch()
}

func (t *Tracer) RestartTheAccumulation() {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.accumulator = NewAccumulator(t)

	t.RootContext.EventOnWireHook = func(e loggeterrapi.IEventInHook, p []byte) {
		t.accumulator.chUp <- p
	}
}

func (t *Tracer) GetUserObject() (uo any) {
	t.mu.Lock()
	defer t.mu.Unlock()

	return t.UserObject
}

// also increments ETC
func (t *Tracer) GetForwardtrace() (tc *trace.Forwardtrace) {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.Forwardtrace.ETC.Inc()

	return t.Forwardtrace
}

func (t *Tracer) SpanBegin(tag string) loggeterrapi.ITracer {
	t.mu.Lock()
	defer t.mu.Unlock()

	tc := t.Forwardtrace

	tc.ETC = append(tc.ETC, 0)
	tc.Tags = append(tc.Tags, tag)

	h := t.RootContext.SpanBeginHook
	if h != nil {
		h(tc)
	}

	return t
}

func (t *Tracer) SpanEnd() {
	t.mu.Lock()
	defer t.mu.Unlock()

	tc := t.Forwardtrace

	h := t.RootContext.SpanEndHook
	if h != nil {
		h(tc)
	}

	tc.ETC = tc.ETC[:len(tc.ETC)-1]
	tc.Tags = tc.Tags[:len(tc.Tags)-1]
}

func (t *Tracer) RaiseCause(first_backtrace_point *trace.BacktracePoint) {
	controlflow.RaiseCause(t.Forwardtrace, first_backtrace_point)
}

func (t *Tracer) NewControlflowEvent(first_backtrace_point *trace.BacktracePoint) *controlflow.Event {
	e := controlflow.NewEvent(t.Forwardtrace, first_backtrace_point)
	return e
}
