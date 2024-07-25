package logger

import "github.com/rs/zerolog"

var _ IEvent = (*Event)(nil)

type Event struct {
	ParentLogger *Logger
	IsInactive   bool   // per event
	TheLevel     string // per event
	TheMsgtag    Msgtag // per event
	TheTitle     string // per event

	subLoggerInitChain *subLoggerInitChain
	zerologEvent       *zerolog.Event
}

type subLoggerInitChain struct {
	zerologContext     *zerolog.Context
	perLoggerSubsystem string
	isInactive         bool // for new sub-logger
}

var _ ILogger = (*Logger)(nil)

type Settings struct {
	PanicOnMisuse  bool
	MsgtagKey      string
	OnSendHook     func(e *Event) (doSend bool)
	ActivationHook func(e *Event) (inactivate bool)
}
type Logger struct {
	Settings     *Settings
	IsInactive   bool   // per whole logger instance
	TheSubsystem string // per logger

	zeroLogger *zerolog.Logger
}
