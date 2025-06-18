package loggeterrzl

import (
	"fmt"
	"time"

	"github.com/rs/zerolog"
	loggeterrapi "github.com/rusriver/nutz/loggeterr-api"
	"github.com/rusriver/nutz/loggeterr-api/trace"
)

type LogEvent struct {
	ParentLogger      *Tracer
	IsInactive        bool   // per event
	TheLevel          string // per event
	TheTitle          string // per event
	theTitleIsDefault bool
	ZerologEvent      *zerolog.Event
}

var _ loggeterrapi.ILogEvent = (*LogEvent)(nil)

func (e *LogEvent) Inactive() loggeterrapi.ILogEvent {
	e.IsInactive = true // per event
	return e
}

func (e *LogEvent) IfActive(f func(ev loggeterrapi.ILogEvent)) loggeterrapi.ILogEvent {
	if !e.IsInactive {
		f(e)
	}
	return e
}

func (e *LogEvent) Caller(skip ...int) loggeterrapi.ILogEvent {
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

func (e *LogEvent) Str(k, v string) loggeterrapi.ILogEvent {
	if e.IsInactive {
		return e // bypass
	}
	e.ZerologEvent.Str(k, v)
	return e
}

func (e *LogEvent) Strf(k, s string, vv ...any) loggeterrapi.ILogEvent {
	if e.IsInactive {
		return e // bypass
	}
	e.ZerologEvent.Str(k, fmt.Sprintf(s, vv...))
	return e
}

func (e *LogEvent) Strs(k string, vv []string) loggeterrapi.ILogEvent {
	if e.IsInactive {
		return e // bypass
	}
	e.ZerologEvent.Strs(k, vv)
	return e
}

func (e *LogEvent) Time(k string, t time.Time) loggeterrapi.ILogEvent {
	if e.IsInactive {
		return e // bypass
	}
	e.ZerologEvent.Time(k, t)
	return e
}

func (e *LogEvent) Int(k string, v int) loggeterrapi.ILogEvent {
	if e.IsInactive {
		return e // bypass
	}
	e.ZerologEvent.Int(k, v)
	return e
}

func (e *LogEvent) Ints(k string, v []int) loggeterrapi.ILogEvent {
	if e.IsInactive {
		return e // bypass
	}
	e.ZerologEvent.Ints(k, v)
	return e
}

func (e *LogEvent) Float64(k string, v float64) loggeterrapi.ILogEvent {
	if e.IsInactive {
		return e // bypass
	}
	e.ZerologEvent.Float64(k, v)
	return e
}

func (e *LogEvent) Floats64(k string, v []float64) loggeterrapi.ILogEvent {
	if e.IsInactive {
		return e // bypass
	}
	e.ZerologEvent.Floats64(k, v)
	return e
}

func (e *LogEvent) Array(k string, v *zerolog.Array) loggeterrapi.ILogEvent {
	if e.IsInactive {
		return e // bypass
	}
	e.ZerologEvent.Array(k, v)
	return e
}

func (e *LogEvent) Dict(k string, v *zerolog.Event) loggeterrapi.ILogEvent {
	if e.IsInactive {
		return e // bypass
	}
	e.ZerologEvent.Dict(k, v)
	return e
}

func (e *LogEvent) SendMsgf(s string, vv ...any) {
	if e.IsInactive {
		return // bypass
	}
	e.ZerologEvent.Str(zerolog.MessageFieldName, fmt.Sprintf(s, vv...))
	e.Send()
}

func (e *LogEvent) SendMsg(s string) {
	if e.IsInactive {
		return // bypass
	}
	e.ZerologEvent.Str(zerolog.MessageFieldName, s)
	e.Send()
}

func (e *LogEvent) Title(s string) loggeterrapi.ILogEvent {
	// this is to protect against losing the title if used twice mistakenly
	switch {
	case len(e.TheTitle) == 0, e.theTitleIsDefault:
		e.TheTitle = s
	default:
		e.TheTitle += " -- " + s
	}
	return e
}

func (e *LogEvent) RawJSON(k string, rj []byte) loggeterrapi.ILogEvent {
	if e.IsInactive {
		return e // bypass
	}
	e.ZerologEvent.RawJSON(k, rj)
	return e
}

func (e *LogEvent) Bytes(k string, bb []byte) loggeterrapi.ILogEvent {
	if e.IsInactive {
		return e // bypass
	}
	e.ZerologEvent.Bytes(k, bb)
	return e
}

func (e *LogEvent) Send() {
	tc := e.ParentLogger.Forwardtrace
	rc := e.ParentLogger.RootContext

	defer func() {
		tc.Tags = tc.Tags[:len(tc.Tags)-1]
	}()

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
	hook := rc.IfSendHook
	doSend := true
	if hook != nil {
		doSend = hook(e)
	}
	if doSend {
		e.ZerologEvent.Send()
	}
}

func (e *LogEvent) GetUserObject() (uo any) {
	e.ParentLogger.mu.Lock()
	defer e.ParentLogger.mu.Unlock()

	return e.ParentLogger.UserObject
}

// doesn't increment ETC (and shouldn't)
func (e *LogEvent) GetForwardtrace() (tc *trace.Forwardtrace) {
	return e.ParentLogger.Forwardtrace
}
