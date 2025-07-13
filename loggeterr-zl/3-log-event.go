package loggeterrzl

import (
	"fmt"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rusriver/nutz/controlflow"
	loggeterrapi "github.com/rusriver/nutz/loggeterr-api"
	"github.com/rusriver/nutz/loggeterr-api/trace"
)

type LogEvent struct {
	ParentLogger      *Tracer
	Tags              []string
	IsInactive        bool
	TheTitle          string
	TheMessage        string
	theTitleIsDefault bool
	ZerologEvent      *zerolog.Event
}

var _ loggeterrapi.ILogEvent = (*LogEvent)(nil)

func (e *LogEvent) Inactive() loggeterrapi.ILogEvent {
	e.IsInactive = true
	return e
}

func (e *LogEvent) IfActive(f func(ev loggeterrapi.ILogEvent)) loggeterrapi.ILogEvent {
	if !e.IsInactive {
		f(e)
	}
	return e
}

func (e *LogEvent) Tag(msgtag string) loggeterrapi.ILogEvent {
	if len(msgtag) > 21 {
		msgtag = msgtag[:21] + "..."
	}
	e.Tags = append(e.Tags, msgtag)
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
	e.TheMessage = fmt.Sprintf(s, vv...)
	e.Send()
}

func (e *LogEvent) SendMsg(x ...any) {
	if e.IsInactive {
		return // bypass
	}
	e.TheMessage = strings.TrimSpace(fmt.Sprintln(x...))
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

func (e *LogEvent) AddControlflowEvent(key string, event *controlflow.Event) loggeterrapi.ILogEvent {
	if e.IsInactive {
		return e // bypass
	}
	bb := event.JSON()
	e.ZerologEvent.RawJSON(key, bb)
	return e
}

func (e *LogEvent) RawStruct(k string, rj []byte) loggeterrapi.ILogEvent {
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
	ft := e.ParentLogger.Forwardtrace
	rc := e.ParentLogger.RootContext

	e.ParentLogger.mu.Lock()
	defer e.ParentLogger.mu.Unlock()
	// this also synchronizes the access to the interceptor

	if e.IsInactive {
		return // bypass
	}

	ft.Tags = append(ft.Tags, e.Tags...)
	defer func() {
		ft.Tags = ft.Tags[:len(ft.Tags)-1]
	}()

	k := rc.ForwardtraceTagsKey
	if len(k) != 0 {
		e.ZerologEvent.Strs(k, ft.Tags)
	}

	k = rc.ForwardtraceETCKey
	if len(k) != 0 {
		ft.ETC.Inc()
		e.ZerologEvent = e.ZerologEvent.Ints(k, ft.ETC)
	}

	if len(e.TheTitle) > 0 {
		e.ZerologEvent.Str("title", e.TheTitle) // gonna be the last, and that's great for readability
	}

	if len(e.TheMessage) > 0 {
		e.ZerologEvent.Str(e.ParentLogger.RootContext.MessageKey, e.TheMessage)
	}

	e.ZerologEvent.Str("tsz", time.Since(e.ParentLogger.RootContext.ZeroTime).String())

	e.ParentLogger.interceptor.e = e // this way we pass it to the EventOnWire() hook
	defer func() { e.ParentLogger.interceptor.e = nil }()

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
