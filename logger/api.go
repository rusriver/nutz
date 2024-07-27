package logger

import (
	"time"

	"github.com/rs/zerolog"
)

type ILogger interface {
	TraceEvent() IEvent
	DebugEvent() IEvent
	ErrorEvent() IEvent
	ErrEvent(error) IEvent
	WarnEvent() IEvent
	InfoEvent() IEvent

	// Can be used if you intend to call ILogger() on it later.
	// By default, all new sub-loggers are active, even if parent was inactivated.
	SubLoggerInitChain() IEvent

	GetZeroLoggerPtr() *zerolog.Logger

	// This is useful to access Settings later after init.
	GetNutzLogger() *Logger

	SetLevel(zerolog.Level) ILogger

	// Attach any user object to the event object, to be used later
	// inside the callback hooks. Intended to pass user-level runtime data to the
	// hook(s). Doesn't add anything to the logline serializations.
	// Normally, the UO must be a pointer, and at the receiving side in the hook
	// you do something along the lines:
	//
	//	myObject, _ := e.ParentLogger.UserObject.(*MyObject)
	//	if myObject !=  nil {
	//		...
	//	}
	AddUO(uo any) ILogger
}

type IEvent interface {
	// Inactivation applies to real events only, all data set ops will be bypassed.
	// However, if applied to sub-logger init chain, nothing is bypassed.
	Inactive() IEvent

	// This is useful to do some work only if the event is activated, e.g. do some data prep
	// or JSON marshalling, and waste no time otherwise.
	//
	//	     logger.InfoEvent().Inactive().Msgtag().IfActive(func(ev IEvent) {
	//				// do some heavy work, then
	//				ev.RawJSON(b)
	//			}).Title("my event").Send()
	IfActive(func(IEvent)) IEvent

	Caller() IEvent
	Str(string, string) IEvent
	Strs(string, []string) IEvent
	Time(string, time.Time) IEvent
	Int(string, int) IEvent
	Array(k string, v *zerolog.Array) (e IEvent)
	Dict(k string, v *zerolog.Event) (e IEvent)

	// Must not be used on sub-logger init chains.
	SendMsgf(string, ...any)

	// Must not be used on sub-logger init chains.
	SendMsg(string)

	// It's like Msg(), but sets the value to "title" field, and also
	// can be used to report to metrics. Must be low-cardinality string.
	// DO NOT put in it any variable strings, e.g. Sprintf()-formatted or
	// containing requestId or any other Id or counters!
	Title(string) IEvent

	// The same as Title().Send()
	// Must not be used on sub-logger init chains.
	SendTitle(string)

	// Some first strings are reported to metrics, be careful to NOT put in them
	// high-cardinality IDs. Doesn't apply to sub-logger init chains, must be set
	// per each event individually.
	// If the ActicationHook is set, it is called, and it can reactivate the event.
	Msgtag(msgtag *Msgtag, ss ...string) IEvent

	// Also must be reported to metrics.
	SubSystem(string) IEvent
	RequestId(string) IEvent
	RawJSON(string, []byte) IEvent
	Bytes(string, []byte) IEvent

	// Must not be used on sub-logger init chains.
	Send()

	// Only works if you created a chain with With()
	ILogger() ILogger
}
