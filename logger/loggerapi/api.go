package loggerapi

import (
	"time"

	"github.com/rs/zerolog"
)

type ILogger interface {
	GetTraceContext() *TraceContext
	GetUserObject() (uo any)

	TraceEvent(msgtag string) IEvent
	DiagnosticEvent(msgtag string) IEvent
	InfoEvent(msgtag string) IEvent
	WarnEvent(msgtag string) IEvent
	ErrorEvent(msgtag string) IEvent
	ErrEvent(err error, msgtag string) IEvent

	// Creates new sub-logger, suitable for another thread.
	// By default, all new sub-loggers are active, even if parent was inactivated.
	NewSubLoggerInitChain(msgtag string) ISubLoggerInitChain
	NSLIC(msgtag string) ISubLoggerInitChain // synonym for NewSubLoggerInitChain()

	// Faster single-thread-only alternative to creating a sub-logger; avoids unnecessary
	// allocations, and calls hooks. Doesn't create another logger, modifies existing one.
	SpanBegin(msgtag string) ILogger
	SpanEnd() // counterpart for SpanBegin()
}

type ISubLoggerInitChain interface {

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
	AddUserObject(uo any) ISubLoggerInitChain

	// Inactivation applies to real events only, all data set ops will be bypassed.
	// However, if applied to sub-logger init chain, nothing is bypassed.
	Inactive() ISubLoggerInitChain

	Caller(skip ...int) ISubLoggerInitChain
	Str(string, string) ISubLoggerInitChain
	Strf(string, string, ...any) ISubLoggerInitChain
	Strs(string, []string) ISubLoggerInitChain
	Time(string, time.Time) ISubLoggerInitChain
	Int(string, int) ISubLoggerInitChain
	Ints(string, []int) ISubLoggerInitChain
	Float64(string, float64) ISubLoggerInitChain
	Floats64(string, []float64) ISubLoggerInitChain
	Array(k string, v *zerolog.Array) ISubLoggerInitChain
	Dict(k string, v *zerolog.Event) ISubLoggerInitChain
	RawJSON(string, []byte) ISubLoggerInitChain
	Bytes(string, []byte) ISubLoggerInitChain

	WithClonedRootContext(optFuncs ...func(rctx *RootContext)) ISubLoggerInitChain
	ILogger() ILogger // (!) Each ILogger must be used in single thread, otherwise you'll get inconsistent tagging
}

type IEvent interface {
	GetTraceContext() *TraceContext
	GetUserObject() (uo any)

	// Inactivation applies to real events only, all data set ops will be bypassed.
	// However, if applied to sub-logger init chain, nothing is bypassed.
	Inactive() IEvent

	// This is useful to do some work only if the event is activated, e.g. do some data prep
	// or JSON marshalling, and waste no time otherwise.
	//
	//	     logger.InfoEvent("1212").Inactive().IfActive(func(ev loggerapi.IEvent) {
	//				// do some heavy work, then
	//				ev.RawJSON(b)
	//			}).Title("my event").Send()
	IfActive(func(IEvent)) IEvent

	Caller(skip ...int) IEvent
	Str(string, string) IEvent
	Strf(string, string, ...any) IEvent
	Strs(string, []string) IEvent
	Time(string, time.Time) IEvent
	Int(string, int) IEvent
	Ints(string, []int) IEvent
	Float64(string, float64) IEvent
	Floats64(string, []float64) IEvent
	Array(k string, v *zerolog.Array) IEvent
	Dict(k string, v *zerolog.Event) IEvent
	RawJSON(string, []byte) IEvent
	Bytes(string, []byte) IEvent

	// It's like Msg(), but sets the value to "title" field, and also
	// can be used to report to metrics. Must be low-cardinality string.
	// DO NOT put in it any variable strings, e.g. Sprintf()-formatted or
	// containing requestId or any other Id or counters!
	// Can't be applied for sub-logger chains, only for events.
	// In typical idiom, must be used BEFORE the Msgtag() - it's specififcally
	// designed to not bypass in inactive events.
	Title(string) IEvent

	// The same as Title().Send()
	SendTitle(string)

	SendMsgf(string, ...any)
	SendMsg(string)
	Send()
}
