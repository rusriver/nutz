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
	SubLoggerInitChain() IEvent
	GetZeroLoggerPtr() *zerolog.Logger
	GetNutzLogger() *Logger
	SetLevel(zerolog.Level) ILogger
	AddUO(uo any) ILogger
}

type IEvent interface {
	Inactive() IEvent
	IfActive(func(IEvent)) IEvent
	Caller() IEvent
	Str(string, string) IEvent
	Strs(string, []string) IEvent
	Time(string, time.Time) IEvent
	Int(string, int) IEvent
	Array(k string, v *zerolog.Array) (e IEvent)
	Dict(k string, v *zerolog.Event) (e IEvent)
	SendMsgf(string, ...any)
	SendMsg(string)
	Title(string) IEvent
	SendTitle(string)
	Msgtag(msgtag *Msgtag, ss ...string) IEvent
	SubSystem(string) IEvent
	RequestId(string) IEvent
	RawJSON(string, []byte) IEvent
	Bytes(string, []byte) IEvent
	Send()
	ILogger() ILogger
}
