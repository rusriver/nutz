package logger

import (
	"time"

	"github.com/rs/zerolog"
)

type ILogger interface {
	Trace() IEvent
	Debug() IEvent
	Error() IEvent
	Err(error) IEvent
	Warn() IEvent
	Info() IEvent
	SubLoggerInitChain() IEvent
	GetZeroLoggerPtr() *zerolog.Logger
	GetNutzLogger() *Logger
	SetLevel(zerolog.Level) ILogger
}

type IEvent interface {
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
