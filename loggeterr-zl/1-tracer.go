package loggeterrzl

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"github.com/rs/zerolog"
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

	ZeroLoggerPtr *zerolog.Logger

	mu          sync.Mutex // only some functions use it
	accumulator *Accumulator
	interceptor *Interceptor
}

var _ loggeterrapi.ITracer = (*Tracer)(nil)
var _ accumapi.ILogsAccumulator = (*Tracer)(nil)

func New(zlPtr *zerolog.Logger, optFuncs ...func(rctx *loggeterrapi.RootContext)) (tracer loggeterrapi.ITracer) {
	rctx := &loggeterrapi.RootContext{}
	// defaults
	rctx.ForwardtraceTagsKey = "FORWARDTRACE-TAGS"
	rctx.ForwardtraceETCKey = "ETC"
	rctx.MessageKey = "MESSAGE"
	rctx.ZeroTime = time.Now()
	for _, f := range optFuncs {
		f(rctx)
	}
	if rctx.MainWriter == nil {
		rctx.MainWriter = os.Stderr
	}
	{
		zerolog.TimestampFieldName = "tsa"
		zerolog.TimeFieldFormat = "20060102-150405"
		zerolog.LevelTraceValue = "PURE-EVENT"
		zerolog.LevelDebugValue = "PURE-EVENT"
		zerolog.LevelInfoValue = "PURE-EVENT"
		zerolog.LevelWarnValue = "PURE-EVENT"
		zerolog.LevelErrorValue = "PURE-EVENT"
		zerolog.LevelFatalValue = "PURE-EVENT"
		zerolog.LevelPanicValue = "PURE-EVENT"
		zerolog.CallerFieldName = "CALLER"
	}
	interceptor := &Interceptor{rctx: rctx}
	if zlPtr == nil {
		zlv := zerolog.New(interceptor).With().Timestamp().Logger()
		zlPtr = &zlv
	}
	tracer = &Tracer{
		RootContext:   rctx,
		ZeroLoggerPtr: zlPtr,
		Forwardtrace: &trace.Forwardtrace{
			ETC: trace.EventTracingCounters{0},
		},
		interceptor: interceptor,
	}
	return tracer
}

func (t *Tracer) EventEvent(msgtag string, reported_e any) loggeterrapi.ILogEvent {
	if len(msgtag) > 21 {
		msgtag = msgtag[:21] + "..."
	}
	e := &LogEvent{
		IsInactive:        t.IsInactive,
		ParentLogger:      t,
		Tags:              []string{msgtag},
		theTitleIsDefault: true,
		ZerologEvent:      t.ZeroLoggerPtr.Info(),
	}
	switch ev := reported_e.(type) {
	case *controlflow.Event:
		bb, _ := json.Marshal(ev)
		e.ZerologEvent.RawJSON("controlflow-event", bb)
		e.TheTitle = "CONTROLFLOW-EVENT"
	case *trace.BacktracePoint:
		bb, _ := json.Marshal(ev)
		e.ZerologEvent.RawJSON("backtracepoint", bb)
		e.TheTitle = "BACKTRACEPOINT"
	case error:
		e.ZerologEvent.Err(ev)
		e.TheTitle = "GO-ERROR"
	default:
		e.ZerologEvent.Str("unknown-event", fmt.Sprintf("%+v", ev))
		e.TheTitle = "UNK"
	}
	activationHook(e)
	return e
}

func (t *Tracer) LogEvent(msgtag string) loggeterrapi.ILogEvent {
	if len(msgtag) > 21 {
		msgtag = msgtag[:21] + "..."
	}

	e := &LogEvent{
		IsInactive:   t.IsInactive,
		ParentLogger: t,
		Tags:         []string{msgtag},
		ZerologEvent: t.ZeroLoggerPtr.Info(),
	}
	activationHook(e)
	return e
}

func activationHook(e *LogEvent) {
	// call the (re-)activation hook, if any
	rc := e.ParentLogger.RootContext
	hook := rc.ActivationHook
	if hook != nil {
		e.IsInactive = hook(e)
	}
}

func (t *Tracer) NewSubTracerInitChain(msgtag string) loggeterrapi.ISubTracerInitChain {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.Forwardtrace.ETC.Inc()

	if len(msgtag) > 32 {
		msgtag = msgtag[:32] + "..."
	}

	zerologContext := t.ZeroLoggerPtr.With()
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

	zlv := t.ZeroLoggerPtr.Output(t.interceptor)
	t.ZeroLoggerPtr = &zlv

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

	if len(tag) > 21 {
		tag = tag[:21] + "..."
	}

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
