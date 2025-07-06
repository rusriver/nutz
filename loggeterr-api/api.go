package loggeterrapi

import (
	"time"

	"github.com/rs/zerolog"
	"github.com/rusriver/nutz/controlflow"
	"github.com/rusriver/nutz/loggeterr-api/trace"
)

type ITracer interface {
	GetForwardtrace() *trace.Forwardtrace // also increments ETC
	GetUserObject() (uo any)

	TraceEvent(msgtag string) ILogEvent
	DiagnosticEvent(msgtag string) ILogEvent
	InfoEvent(msgtag string) ILogEvent
	WarnEvent(msgtag string) ILogEvent
	ErrorEvent(msgtag string) ILogEvent
	SmartErrEvent(err any, msgtag string) ILogEvent // can accept error, *controlflow.Event, *trace.BacktracePoint
	LogEvent(msgtag string) ILogEvent               // without legacy level, pure event; may default to INFO

	// Creates new sub-logger, suitable for another thread.
	// By default, all new sub-loggers are active, even if parent was inactivated.
	NewSubTracerInitChain(msgtag string) ISubTracerInitChain
	NSTIC(msgtag string) ISubTracerInitChain // synonym for NewSubLoggerInitChain()

	// Faster single-thread-only alternative to creating a sub-logger; avoids unnecessary
	// allocations, and calls hooks. Doesn't create another logger, modifies existing one.
	SpanBegin(msgtag string) ITracer
	SpanEnd() // counterpart for SpanBegin()

	// example use: tracer.RaiseCause(trace.BacktracePointf("123123", "the reason"))
	RaiseCause(first_backtrace_point *trace.BacktracePoint)
	NewControlflowEvent(first_backtrace_point *trace.BacktracePoint) *controlflow.Event
}

type ISubTracerInitChain interface {

	// Attach any user object to the event object, to be used later
	// inside the callback hooks. Intended to pass user-level runtime data to the
	// hook(s). Doesn't add anything to the logline serializations.
	// Normally, the UO must be a pointer, and at the receiving side in the hook
	// you do something along the lines:
	//
	//	myObject, _ := e.GetUserObject().(*MyObject)
	//	if myObject !=  nil {
	//		...
	//	}
	SetUserObject(uo any) ISubTracerInitChain

	// Can be used to use custom/external trace
	SetForwardtrace(ft *trace.Forwardtrace) ISubTracerInitChain

	// Inactivation applies to real events only, all data set ops will be bypassed.
	// However, if applied to sub-logger init chain, nothing is bypassed.
	Inactive() ISubTracerInitChain

	Caller(skip ...int) ISubTracerInitChain
	Str(string, string) ISubTracerInitChain
	Strf(string, string, ...any) ISubTracerInitChain
	Strs(string, []string) ISubTracerInitChain
	Time(string, time.Time) ISubTracerInitChain
	Int(string, int) ISubTracerInitChain
	Ints(string, []int) ISubTracerInitChain
	Float64(string, float64) ISubTracerInitChain
	Floats64(string, []float64) ISubTracerInitChain
	Array(k string, v *zerolog.Array) ISubTracerInitChain
	Dict(k string, v *zerolog.Event) ISubTracerInitChain
	RawJSON(string, []byte) ISubTracerInitChain
	Bytes(string, []byte) ISubTracerInitChain

	WithClonedRootContext(optFuncs ...func(rctx *RootContext)) ISubTracerInitChain
	ITracer() ITracer // (!) Each ILogger must be used in single thread, otherwise you'll get inconsistent tagging
}

type ILogEvent interface {
	GetForwardtrace() *trace.Forwardtrace // doesn't increment ETC, used in callbacks
	GetUserObject() (uo any)

	// Inactivation applies to real events only, all data set ops will be bypassed.
	// However, if applied to sub-logger init chain, nothing is bypassed.
	Inactive() ILogEvent

	// This is useful to do some work only if the event is activated, e.g. do some data prep
	// or JSON marshalling, and waste no time otherwise.
	//
	//	     logger.InfoEvent("1212").Inactive().IfActive(func(ev loggerapi.IEvent) {
	//				// do some heavy work, then
	//				ev.RawJSON(b)
	//			}).Title("my event").Send()
	IfActive(func(ILogEvent)) ILogEvent

	Caller(skip ...int) ILogEvent
	Str(string, string) ILogEvent
	Strf(string, string, ...any) ILogEvent
	Strs(string, []string) ILogEvent
	Time(string, time.Time) ILogEvent
	Int(string, int) ILogEvent
	Ints(string, []int) ILogEvent
	Float64(string, float64) ILogEvent
	Floats64(string, []float64) ILogEvent
	Array(k string, v *zerolog.Array) ILogEvent
	Dict(k string, v *zerolog.Event) ILogEvent
	RawJSON(string, []byte) ILogEvent
	Bytes(string, []byte) ILogEvent

	// It's like Msg(), but sets the value to "title" field, and also
	// can be used to report to metrics. Must be low-cardinality string.
	// DO NOT put in it any variable strings, e.g. Sprintf()-formatted or
	// containing requestId or any other Id or counters!
	// Can't be applied for sub-logger chains, only for events.
	// In typical idiom, must be used BEFORE the Msgtag() - it's specififcally
	// designed to not bypass in inactive events.
	Title(string) ILogEvent

	AddControlflowEvent(key string, event *controlflow.Event) ILogEvent

	SendMsgf(string, ...any)
	SendMsg(...any)
	Send()
}
