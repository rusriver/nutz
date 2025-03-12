package logger

import (
	"fmt"
	"time"

	"github.com/rs/zerolog"
	"github.com/rusriver/nutz/logger/loggerapi"
)

type Event struct {
	ParentLogger      *Logger
	IsInactive        bool              // per event
	TheLevel          string            // per event
	TheMsgtags        loggerapi.Msgtags // per event
	TheTitle          string            // per event
	theTitleIsDefault bool
	ZerologEvent      *zerolog.Event
}

var _ loggerapi.IEvent = (*Event)(nil)

func (e *Event) Inactive() loggerapi.IEvent {
	e.IsInactive = true // per event
	return e
}

func (e *Event) IfActive(f func(ev loggerapi.IEvent)) loggerapi.IEvent {
	if !e.IsInactive {
		f(e)
	}
	return e
}

func (e *Event) Caller(skip ...int) loggerapi.IEvent {
	skipCount := 4
	if len(skip) > 0 {
		skipCount = skip[0]
	}
	if e.IsInactive {
		return e // bypass
	}
	e.ZerologEvent.Caller(skipCount)
	return e
}

func (e *Event) Str(k, v string) loggerapi.IEvent {
	if e.IsInactive {
		return e // bypass
	}
	e.ZerologEvent.Str(k, v)
	return e
}

func (e *Event) Strf(k, s string, vv ...any) loggerapi.IEvent {
	if e.IsInactive {
		return e // bypass
	}
	e.ZerologEvent.Str(k, fmt.Sprintf(s, vv...))
	return e
}

func (e *Event) Strs(k string, vv []string) loggerapi.IEvent {
	if e.IsInactive {
		return e // bypass
	}
	e.ZerologEvent.Strs(k, vv)
	return e
}

func (e *Event) Time(k string, t time.Time) loggerapi.IEvent {
	if e.IsInactive {
		return e // bypass
	}
	e.ZerologEvent.Time(k, t)
	return e
}

func (e *Event) Int(k string, v int) loggerapi.IEvent {
	if e.IsInactive {
		return e // bypass
	}
	e.ZerologEvent.Int(k, v)
	return e
}

func (e *Event) Ints(k string, v []int) loggerapi.IEvent {
	if e.IsInactive {
		return e // bypass
	}
	e.ZerologEvent.Ints(k, v)
	return e
}

func (e *Event) Float64(k string, v float64) loggerapi.IEvent {
	if e.IsInactive {
		return e // bypass
	}
	e.ZerologEvent.Float64(k, v)
	return e
}

func (e *Event) Floats64(k string, v []float64) loggerapi.IEvent {
	if e.IsInactive {
		return e // bypass
	}
	e.ZerologEvent.Floats64(k, v)
	return e
}

func (e *Event) Array(k string, v *zerolog.Array) loggerapi.IEvent {
	if e.IsInactive {
		return e // bypass
	}
	e.ZerologEvent.Array(k, v)
	return e
}

func (e *Event) Dict(k string, v *zerolog.Event) loggerapi.IEvent {
	if e.IsInactive {
		return e // bypass
	}
	e.ZerologEvent.Dict(k, v)
	return e
}

func (e *Event) SendMsgf(s string, vv ...any) {
	if e.IsInactive {
		return // bypass
	}
	e.ZerologEvent.Str(zerolog.MessageFieldName, fmt.Sprintf(s, vv...))
	e.Send()
}

func (e *Event) SendMsg(s string) {
	if e.IsInactive {
		return // bypass
	}
	e.ZerologEvent.Str(zerolog.MessageFieldName, s)
	e.Send()
}

func (e *Event) Title(s string) loggerapi.IEvent {
	// this is to protect against losing the title if used twice mistakenly
	switch {
	case len(e.TheTitle) == 0, e.theTitleIsDefault:
		e.TheTitle = s
	default:
		e.TheTitle += " -- " + s
	}
	return e
}

func (e *Event) SendTitle(s string) {
	if e.IsInactive {
		return // bypass
	}
	e.Title(s).Send()
}

func (e *Event) RawJSON(k string, rj []byte) loggerapi.IEvent {
	if e.IsInactive {
		return e // bypass
	}
	e.ZerologEvent.RawJSON(k, rj)
	return e
}

func (e *Event) Bytes(k string, bb []byte) loggerapi.IEvent {
	if e.IsInactive {
		return e // bypass
	}
	e.ZerologEvent.Bytes(k, bb)
	return e
}

func (e *Event) Send() {
	if e.IsInactive {
		return // bypass
	}

	e.ParentLogger.mu.Lock()
	defer e.ParentLogger.mu.Unlock()
	// this also synchronizes access to the interceptor

	e.ParentLogger.interceptor.e = e // this way we pass it to the EventOnWire() hook
	defer func() { e.ParentLogger.interceptor.e = nil }()

	if len(e.TheTitle) > 0 {
		e.ZerologEvent.Str("title", e.TheTitle) // gonna be the last, and that's great for readability
	}
	hook := e.ParentLogger.RootContext.IfSendHook
	doSend := true
	if hook != nil {
		doSend = hook(e)
	}
	if doSend {
		e.ZerologEvent.Send()
	}
}

func (e *Event) GetUserObject() (uo any) {
	e.ParentLogger.mu.Lock()
	defer e.ParentLogger.mu.Unlock()

	return e.ParentLogger.UserObject
}

func (e *Event) GetTraceContext() (tc *loggerapi.TraceContext) {
	return e.ParentLogger.TraceContext
}
