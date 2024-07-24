package logger

import (
	"fmt"
	"time"

	"github.com/rs/zerolog"
)

var _ IEvent = (*Event)(nil)

type Event struct {
	ParentLogger *Logger
	TheLevel     string // per event
	TheMsgtag    Msgtag // per event
	TheTitle     string // per event

	subLoggerInitChain *subLoggerInitChain
	zerologEvent       *zerolog.Event
}

type subLoggerInitChain struct {
	zerologContext     *zerolog.Context
	perLoggerSubsystem string
}

func (e *Event) Caller() IEvent {
	if e.subLoggerInitChain != nil {
		sub := e.subLoggerInitChain.zerologContext.Caller()
		e.subLoggerInitChain.zerologContext = &sub
	} else {
		e.zerologEvent.Caller()
	}
	return e
}

func (e *Event) Str(k, v string) IEvent {
	if e.subLoggerInitChain != nil {
		sub := e.subLoggerInitChain.zerologContext.Str(k, v)
		e.subLoggerInitChain.zerologContext = &sub
	} else {
		e.zerologEvent.Str(k, v)
	}
	return e
}

func (e *Event) Strs(k string, vv []string) IEvent {
	if e.subLoggerInitChain != nil {
		sub := e.subLoggerInitChain.zerologContext.Strs(k, vv)
		e.subLoggerInitChain.zerologContext = &sub
	} else {
		e.zerologEvent.Strs(k, vv)
	}
	return e
}

func (e *Event) Time(k string, t time.Time) IEvent {
	if e.subLoggerInitChain != nil {
		sub := e.subLoggerInitChain.zerologContext.Time(k, t)
		e.subLoggerInitChain.zerologContext = &sub
	} else {
		e.zerologEvent.Time(k, t)
	}
	return e
}

func (e *Event) Int(k string, v int) IEvent {
	if e.subLoggerInitChain != nil {
		sub := e.subLoggerInitChain.zerologContext.Int(k, v)
		e.subLoggerInitChain.zerologContext = &sub
	} else {
		e.zerologEvent.Int(k, v)
	}
	return e
}

func (e *Event) Array(k string, v *zerolog.Array) IEvent {
	if e.subLoggerInitChain != nil {
		sub := e.subLoggerInitChain.zerologContext.Array(k, v)
		e.subLoggerInitChain.zerologContext = &sub
	} else {
		e.zerologEvent.Array(k, v)
	}
	return e
}

func (e *Event) Dict(k string, v *zerolog.Event) IEvent {
	if e.subLoggerInitChain != nil {
		sub := e.subLoggerInitChain.zerologContext.Dict(k, v)
		e.subLoggerInitChain.zerologContext = &sub
	} else {
		e.zerologEvent.Dict(k, v)
	}
	return e
}

// Must not be used on sub-logger init chains.
func (e *Event) SendMsgf(s string, vv ...any) {
	e.zerologEvent.Str(zerolog.MessageFieldName, fmt.Sprintf(s, vv...))
	e.Send()
	return
}

// Must not be used on sub-logger init chains.
func (e *Event) SendMsg(s string) {
	e.zerologEvent.Str(zerolog.MessageFieldName, s)
	e.Send()
	return
}

// It's like Msg(), but sets the value to "title" field, and also
// can be used to report to metrics. Must be low-cardinality string.
// DO NOT put in it any variable strings, e.g. Sprintf()-formatted or
// containing requestId or any other Id or counters!
func (e *Event) Title(s string) IEvent {
	if e.subLoggerInitChain != nil {
		sub := e.subLoggerInitChain.zerologContext.Str("title", s)
		e.subLoggerInitChain.zerologContext = &sub
	} else {
		e.zerologEvent.Str("title", s)
		e.TheTitle = s
	}
	return e
}

// The same as Title().Send()
// Must not be used on sub-logger init chains.
func (e *Event) SendTitle(s string) {
	e.Title(s).Send()
	return
}

// Some first strings are reported to metrics, be careful to NOT put in them
// high-cardinality IDs. Doesn't apply to sub-logger init chains, must be set
// per each event individually.
func (e *Event) Msgtag(msgtag *Msgtag, ss ...string) IEvent {
	if msgtag == nil {
		msgtag = &Msgtag{}
	}
	msgtag = msgtag.With(ss...)
	e.TheMsgtag = *msgtag
	if len(*msgtag) > 0 {
		e.zerologEvent.Strs(e.ParentLogger.Settings.MsgtagKey, *msgtag)
	}
	return e
}

// Also must be reported to metrics.
func (e *Event) SubSystem(s string) IEvent {
	if e.subLoggerInitChain != nil {
		sub := e.subLoggerInitChain.zerologContext.Str("subsystem", s)
		e.subLoggerInitChain.zerologContext = &sub
		e.subLoggerInitChain.perLoggerSubsystem = s
	} else {
		e.zerologEvent.Str("subsystem", s)
	}
	return e
}

func (e *Event) RequestId(s string) IEvent {
	if e.subLoggerInitChain != nil {
		sub := e.subLoggerInitChain.zerologContext.Str("requestId", s)
		e.subLoggerInitChain.zerologContext = &sub
	} else {
		e.zerologEvent.Str("requestId", s)
	}
	return e
}

func (e *Event) RawJSON(k string, rj []byte) IEvent {
	if e.subLoggerInitChain != nil {
		sub := e.subLoggerInitChain.zerologContext.RawJSON(k, rj)
		e.subLoggerInitChain.zerologContext = &sub
	} else {
		e.zerologEvent.RawJSON(k, rj)
	}
	return e
}

func (e *Event) Bytes(k string, bb []byte) IEvent {
	if e.subLoggerInitChain != nil {
		sub := e.subLoggerInitChain.zerologContext.Bytes(k, bb)
		e.subLoggerInitChain.zerologContext = &sub
	} else {
		e.zerologEvent.Bytes(k, bb)
	}
	return e
}

// Must not be used on sub-logger init chains.
func (e *Event) Send() {
	if e.subLoggerInitChain != nil {
		if e.ParentLogger.Settings.PanicOnMisuse {
			panic("aeb65f98-e2ff-489d-ac0f-b8416c3c1e78")
		} else {
			return
		}
	}
	hook := e.ParentLogger.Settings.OnSendHook
	doSend := true
	if hook != nil {
		doSend = hook(e)
	}
	if doSend {
		e.zerologEvent.Send()
	}
	return
}

// Only works if you created a chain with With()
func (e *Event) ILogger() ILogger {
	if e.subLoggerInitChain != nil {
		zeroLogger := e.subLoggerInitChain.zerologContext.Logger()
		logger := &Logger{
			Settings:     e.ParentLogger.Settings,
			TheSubsystem: e.subLoggerInitChain.perLoggerSubsystem,
			zeroLogger:   &zeroLogger,
		}
		return logger
	} else {
		if e.ParentLogger.Settings.PanicOnMisuse {
			panic("5a2d0286-d0bd-4adc-a2ce-da0f68f3b0fb")
		} else {
			return nil
		}
	}
}
