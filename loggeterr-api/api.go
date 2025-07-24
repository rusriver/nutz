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

	/*
		Levels:
			trace, diagnostic, informational, warning - deprecated.
			error - discouraged to use, and found to be ambiguous in most real-life non-trivial apps
			fatal - outlawed long ago

		See also https://dave.cheney.net/2015/11/05/lets-talk-about-logging

		Long story short, what if you really want semantic levels? JUST USE FORWARD TRACE TAGS.
		Consider this point of view: what if these levels historically were just a forerunners
		for forward traces.

		See also the Tag() method, it can help a lot.
	*/

	EventEvent(msgtag string, e any) ILogEvent // can accept error, *controlflow.Event, *trace.BacktracePoint // max tag length is 21
	LogEvent(msgtag string) ILogEvent          // pure log event, without legacy semantic level // max tag length is 21

	// Creates new sub-logger, suitable for another thread.
	// By default, all new sub-loggers are active, even if parent was inactivated.
	NewSubTracerInitChain(msgtag string) ISubTracerInitChain // max tag length is 32
	NSTIC(msgtag string) ISubTracerInitChain                 // synonym for NewSubLoggerInitChain()

	// Faster single-thread-only alternative to creating a sub-logger; avoids unnecessary
	// allocations, and calls hooks. Doesn't create another logger, modifies existing one.
	SpanBegin(msgtag string) ITracer // max tag length is 21
	SpanEnd()                        // counterpart for SpanBegin()

	Reset() // resets the tracer's anti-concurrency protection, if previous event wasn't sent

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
	RawStruct(string, []byte) ISubTracerInitChain // sets data as is, the structural correctness is responsibility of user
	Bytes(string, []byte) ISubTracerInitChain     // sets data as valid literal

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
	//	     logger.LogEvent("1212").Inactive().IfActive(func(ev loggerapi.IEvent) {
	//				// do some heavy work, then
	//				ev.RawJSON(b)
	//			}).Title("my event").Send()
	IfActive(func(ILogEvent)) ILogEvent

	Tag(msgtag string) ILogEvent // adds extra tag to forward trace of event // max length is 21

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
	RawStruct(string, []byte) ILogEvent // sets data as is, the structural correctness is responsibility of user
	Bytes(string, []byte) ILogEvent     // sets data as valid literal

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
