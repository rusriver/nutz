package logger

import (
	"fmt"
	"time"

	"github.com/rs/zerolog"
)

func (e *Event) Inactive() IEvent {
	if e.subLoggerInitChain != nil {
		e.subLoggerInitChain.isInactive = true // for new sub-logger init chain
	} else {
		e.IsInactive = true // per event
	}
	return e
}

func (e *Event) IfActive(f func(ev IEvent)) IEvent {
	if !e.IsInactive {
		f(e)
	}
	return e
}

func (e *Event) Caller(skip ...int) IEvent {
	skipCount := 4
	if len(skip) > 0 {
		skipCount = skip[0]
	}
	if e.subLoggerInitChain != nil {
		sub := e.subLoggerInitChain.zerologContext.CallerWithSkipFrameCount(skipCount)
		e.subLoggerInitChain.zerologContext = &sub
	} else {
		if e.IsInactive {
			return e // bypass
		}
		e.zerologEvent.Caller(skipCount)
	}
	return e
}

func (e *Event) Str(k, v string) IEvent {
	if e.subLoggerInitChain != nil {
		sub := e.subLoggerInitChain.zerologContext.Str(k, v)
		e.subLoggerInitChain.zerologContext = &sub
	} else {
		if e.IsInactive {
			return e // bypass
		}
		e.zerologEvent.Str(k, v)
	}
	return e
}

func (e *Event) Strf(k, s string, vv ...any) IEvent {
	if e.subLoggerInitChain != nil {
		sub := e.subLoggerInitChain.zerologContext.Str(k, fmt.Sprintf(s, vv...))
		e.subLoggerInitChain.zerologContext = &sub
	} else {
		if e.IsInactive {
			return e // bypass
		}
		e.zerologEvent.Str(k, fmt.Sprintf(s, vv...))
	}
	return e
}

func (e *Event) Strs(k string, vv []string) IEvent {
	if e.subLoggerInitChain != nil {
		sub := e.subLoggerInitChain.zerologContext.Strs(k, vv)
		e.subLoggerInitChain.zerologContext = &sub
	} else {
		if e.IsInactive {
			return e // bypass
		}
		e.zerologEvent.Strs(k, vv)
	}
	return e
}

func (e *Event) Time(k string, t time.Time) IEvent {
	if e.subLoggerInitChain != nil {
		sub := e.subLoggerInitChain.zerologContext.Time(k, t)
		e.subLoggerInitChain.zerologContext = &sub
	} else {
		if e.IsInactive {
			return e // bypass
		}
		e.zerologEvent.Time(k, t)
	}
	return e
}

func (e *Event) Int(k string, v int) IEvent {
	if e.subLoggerInitChain != nil {
		sub := e.subLoggerInitChain.zerologContext.Int(k, v)
		e.subLoggerInitChain.zerologContext = &sub
	} else {
		if e.IsInactive {
			return e // bypass
		}
		e.zerologEvent.Int(k, v)
	}
	return e
}

func (e *Event) Float64(k string, v float64) IEvent {
	if e.subLoggerInitChain != nil {
		sub := e.subLoggerInitChain.zerologContext.Float64(k, v)
		e.subLoggerInitChain.zerologContext = &sub
	} else {
		if e.IsInactive {
			return e // bypass
		}
		e.zerologEvent.Float64(k, v)
	}
	return e
}

func (e *Event) Array(k string, v *zerolog.Array) IEvent {
	if e.subLoggerInitChain != nil {
		sub := e.subLoggerInitChain.zerologContext.Array(k, v)
		e.subLoggerInitChain.zerologContext = &sub
	} else {
		if e.IsInactive {
			return e // bypass
		}
		e.zerologEvent.Array(k, v)
	}
	return e
}

func (e *Event) Dict(k string, v *zerolog.Event) IEvent {
	if e.subLoggerInitChain != nil {
		sub := e.subLoggerInitChain.zerologContext.Dict(k, v)
		e.subLoggerInitChain.zerologContext = &sub
	} else {
		if e.IsInactive {
			return e // bypass
		}
		e.zerologEvent.Dict(k, v)
	}
	return e
}

func (e *Event) SendMsgf(s string, vv ...any) {
	if e.IsInactive {
		return // bypass
	}
	e.zerologEvent.Str(zerolog.MessageFieldName, fmt.Sprintf(s, vv...))
	e.Send()
	return
}

func (e *Event) SendMsg(s string) {
	if e.IsInactive {
		return // bypass
	}
	e.zerologEvent.Str(zerolog.MessageFieldName, s)
	e.Send()
	return
}

func (e *Event) Title(s string) IEvent {
	if e.subLoggerInitChain != nil {
		if e.ParentLogger.Settings.PanicOnMisuse {
			panic("282dd247-6dbd-410f-989b-59ef2daa67ff")
		} else {
			return e
		}
	} else {
		// this is to protect against losing the title if used twice mistakenly
		if len(e.TheTitle) == 0 {
			e.TheTitle = s
		} else {
			e.TheTitle += " " + s
		}
	}
	return e
}

func (e *Event) SendTitle(s string) {
	if e.IsInactive {
		return // bypass
	}
	e.Title(s).Send()
	return
}

func (e *Event) Msgtag(msgtag *Msgtag, ss ...string) IEvent {
	// set msgtag to our context first
	if msgtag == nil {
		msgtag = &Msgtag{}
	}
	msgtag = msgtag.With(ss...)
	e.TheMsgtag = *msgtag

	// call the (re-)activation hook, if any
	hook := e.ParentLogger.Settings.ActivationHook
	if hook != nil {
		e.IsInactive = hook(e)
	}

	// then handle the rest as usual
	if e.IsInactive {
		return e // bypass
	}
	if len(*msgtag) > 0 {
		e.zerologEvent.Strs(e.ParentLogger.Settings.MsgtagKey, *msgtag)
	}
	return e
}

func (e *Event) SubSystem(s string) IEvent {
	if e.subLoggerInitChain != nil {
		sub := e.subLoggerInitChain.zerologContext.Str("subsystem", s)
		e.subLoggerInitChain.zerologContext = &sub
		e.subLoggerInitChain.perLoggerSubsystem = s
	} else {
		if e.IsInactive {
			return e // bypass
		}
		e.zerologEvent.Str("subsystem", s)
	}
	return e
}

func (e *Event) RequestId(s string) IEvent {
	if e.subLoggerInitChain != nil {
		sub := e.subLoggerInitChain.zerologContext.Str("requestId", s)
		e.subLoggerInitChain.zerologContext = &sub
	} else {
		if e.IsInactive {
			return e // bypass
		}
		e.zerologEvent.Str("requestId", s)
	}
	return e
}

func (e *Event) RawJSON(k string, rj []byte) IEvent {
	if e.subLoggerInitChain != nil {
		sub := e.subLoggerInitChain.zerologContext.RawJSON(k, rj)
		e.subLoggerInitChain.zerologContext = &sub
	} else {
		if e.IsInactive {
			return e // bypass
		}
		e.zerologEvent.RawJSON(k, rj)
	}
	return e
}

func (e *Event) Bytes(k string, bb []byte) IEvent {
	if e.subLoggerInitChain != nil {
		sub := e.subLoggerInitChain.zerologContext.Bytes(k, bb)
		e.subLoggerInitChain.zerologContext = &sub
	} else {
		if e.IsInactive {
			return e // bypass
		}
		e.zerologEvent.Bytes(k, bb)
	}
	return e
}

func (e *Event) Send() {
	if e.subLoggerInitChain != nil {
		if e.ParentLogger.Settings.PanicOnMisuse {
			panic("aeb65f98-e2ff-489d-ac0f-b8416c3c1e78")
		} else {
			return
		}
	}
	if e.IsInactive {
		return // bypass
	}
	if len(e.TheTitle) > 0 {
		e.zerologEvent.Str("title", e.TheTitle) // gonna be the last, and that's great for readability
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

func (e *Event) ILogger() ILogger {
	if e.subLoggerInitChain != nil {
		zeroLogger := e.subLoggerInitChain.zerologContext.Logger()
		logger := &Logger{
			Settings:     e.ParentLogger.Settings,
			IsInactive:   e.subLoggerInitChain.isInactive, // inherit from the init chain
			TheSubsystem: e.subLoggerInitChain.perLoggerSubsystem,
			UserObject:   e.ParentLogger.UserObject,
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
