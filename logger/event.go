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

	subLoggerInitChain *subLoggerInitChain
	ZerologEvent       *zerolog.Event
}

var _ loggerapi.IEvent = (*Event)(nil)

func (e *Event) Inactive() loggerapi.IEvent {
	if e.subLoggerInitChain != nil {
		e.subLoggerInitChain.isInactive = true // for new sub-logger init chain
	} else {
		e.IsInactive = true // per event
	}
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
	if e.subLoggerInitChain != nil {
		sub := e.subLoggerInitChain.zerologContext.CallerWithSkipFrameCount(skipCount)
		e.subLoggerInitChain.zerologContext = &sub
	} else {
		if e.IsInactive {
			return e // bypass
		}
		e.ZerologEvent.Caller(skipCount)
	}
	return e
}

func (e *Event) Str(k, v string) loggerapi.IEvent {
	if e.subLoggerInitChain != nil {
		sub := e.subLoggerInitChain.zerologContext.Str(k, v)
		e.subLoggerInitChain.zerologContext = &sub
	} else {
		if e.IsInactive {
			return e // bypass
		}
		e.ZerologEvent.Str(k, v)
	}
	return e
}

func (e *Event) Strf(k, s string, vv ...any) loggerapi.IEvent {
	if e.subLoggerInitChain != nil {
		sub := e.subLoggerInitChain.zerologContext.Str(k, fmt.Sprintf(s, vv...))
		e.subLoggerInitChain.zerologContext = &sub
	} else {
		if e.IsInactive {
			return e // bypass
		}
		e.ZerologEvent.Str(k, fmt.Sprintf(s, vv...))
	}
	return e
}

func (e *Event) Strs(k string, vv []string) loggerapi.IEvent {
	if e.subLoggerInitChain != nil {
		sub := e.subLoggerInitChain.zerologContext.Strs(k, vv)
		e.subLoggerInitChain.zerologContext = &sub
	} else {
		if e.IsInactive {
			return e // bypass
		}
		e.ZerologEvent.Strs(k, vv)
	}
	return e
}

func (e *Event) Time(k string, t time.Time) loggerapi.IEvent {
	if e.subLoggerInitChain != nil {
		sub := e.subLoggerInitChain.zerologContext.Time(k, t)
		e.subLoggerInitChain.zerologContext = &sub
	} else {
		if e.IsInactive {
			return e // bypass
		}
		e.ZerologEvent.Time(k, t)
	}
	return e
}

func (e *Event) Int(k string, v int) loggerapi.IEvent {
	if e.subLoggerInitChain != nil {
		sub := e.subLoggerInitChain.zerologContext.Int(k, v)
		e.subLoggerInitChain.zerologContext = &sub
	} else {
		if e.IsInactive {
			return e // bypass
		}
		e.ZerologEvent.Int(k, v)
	}
	return e
}

func (e *Event) Ints(k string, v []int) loggerapi.IEvent {
	if e.subLoggerInitChain != nil {
		sub := e.subLoggerInitChain.zerologContext.Ints(k, v)
		e.subLoggerInitChain.zerologContext = &sub
	} else {
		if e.IsInactive {
			return e // bypass
		}
		e.ZerologEvent.Ints(k, v)
	}
	return e
}

func (e *Event) Float64(k string, v float64) loggerapi.IEvent {
	if e.subLoggerInitChain != nil {
		sub := e.subLoggerInitChain.zerologContext.Float64(k, v)
		e.subLoggerInitChain.zerologContext = &sub
	} else {
		if e.IsInactive {
			return e // bypass
		}
		e.ZerologEvent.Float64(k, v)
	}
	return e
}

func (e *Event) Floats64(k string, v []float64) loggerapi.IEvent {
	if e.subLoggerInitChain != nil {
		sub := e.subLoggerInitChain.zerologContext.Floats64(k, v)
		e.subLoggerInitChain.zerologContext = &sub
	} else {
		if e.IsInactive {
			return e // bypass
		}
		e.ZerologEvent.Floats64(k, v)
	}
	return e
}

func (e *Event) Array(k string, v *zerolog.Array) loggerapi.IEvent {
	if e.subLoggerInitChain != nil {
		sub := e.subLoggerInitChain.zerologContext.Array(k, v)
		e.subLoggerInitChain.zerologContext = &sub
	} else {
		if e.IsInactive {
			return e // bypass
		}
		e.ZerologEvent.Array(k, v)
	}
	return e
}

func (e *Event) Dict(k string, v *zerolog.Event) loggerapi.IEvent {
	if e.subLoggerInitChain != nil {
		sub := e.subLoggerInitChain.zerologContext.Dict(k, v)
		e.subLoggerInitChain.zerologContext = &sub
	} else {
		if e.IsInactive {
			return e // bypass
		}
		e.ZerologEvent.Dict(k, v)
	}
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
	if e.subLoggerInitChain != nil {
		if e.ParentLogger.RootContext.PanicOnMisuse {
			panic("282dd247-6dbd-410f-989b-59ef2daa67ff")
		} else {
			return e
		}
	} else {
		// this is to protect against losing the title if used twice mistakenly
		switch {
		case len(e.TheTitle) == 0, e.theTitleIsDefault:
			e.TheTitle = s
		default:
			e.TheTitle += " -- " + s
		}
	}
	return e
}

func (e *Event) SendTitle(s string) {
	if e.IsInactive {
		return // bypass
	}
	e.Title(s).Send()
}

func (e *Event) SubSystem(s string) loggerapi.IEvent {
	if e.subLoggerInitChain != nil {
		sub := e.subLoggerInitChain.zerologContext.Str("subsystem", s)
		e.subLoggerInitChain.zerologContext = &sub
		e.subLoggerInitChain.perLoggerSubsystem = s
	} else {
		if e.IsInactive {
			return e // bypass
		}
		e.ZerologEvent.Str("subsystem", s)
	}
	return e
}

func (e *Event) RequestId(s string) loggerapi.IEvent {
	if e.subLoggerInitChain != nil {
		sub := e.subLoggerInitChain.zerologContext.Str("requestId", s)
		e.subLoggerInitChain.zerologContext = &sub
	} else {
		if e.IsInactive {
			return e // bypass
		}
		e.ZerologEvent.Str("requestId", s)
	}
	return e
}

func (e *Event) RawJSON(k string, rj []byte) loggerapi.IEvent {
	if e.subLoggerInitChain != nil {
		sub := e.subLoggerInitChain.zerologContext.RawJSON(k, rj)
		e.subLoggerInitChain.zerologContext = &sub
	} else {
		if e.IsInactive {
			return e // bypass
		}
		e.ZerologEvent.RawJSON(k, rj)
	}
	return e
}

func (e *Event) Bytes(k string, bb []byte) loggerapi.IEvent {
	if e.subLoggerInitChain != nil {
		sub := e.subLoggerInitChain.zerologContext.Bytes(k, bb)
		e.subLoggerInitChain.zerologContext = &sub
	} else {
		if e.IsInactive {
			return e // bypass
		}
		e.ZerologEvent.Bytes(k, bb)
	}
	return e
}

func (e *Event) Send() {
	e.ParentLogger.mu.Lock()
	defer e.ParentLogger.mu.Unlock()

	if e.subLoggerInitChain != nil {
		if e.ParentLogger.RootContext.PanicOnMisuse {
			panic("aeb65f98-e2ff-489d-ac0f-b8416c3c1e78")
		} else {
			return
		}
	}
	if e.IsInactive {
		return // bypass
	}
	if len(e.TheTitle) > 0 {
		e.ZerologEvent.Str("title", e.TheTitle) // gonna be the last, and that's great for readability
	}
	hook := e.ParentLogger.RootContext.OnSendHook
	doSend := true
	if hook != nil {
		doSend = hook(e)
	}
	if doSend {
		e.ZerologEvent.Send()
	}
}

func (e *Event) ILogger() loggerapi.ILogger {
	e.ParentLogger.mu.Lock()
	defer e.ParentLogger.mu.Unlock()

	if e.subLoggerInitChain != nil {
		zeroLogger := e.subLoggerInitChain.zerologContext.Logger()
		logger := &Logger{
			RootContext:  e.ParentLogger.RootContext,
			IsInactive:   e.subLoggerInitChain.isInactive, // inherit from the init chain
			TheSubsystem: e.subLoggerInitChain.perLoggerSubsystem,
			UserObject:   e.ParentLogger.UserObject,
			ZeroLogger:   &zeroLogger,
			ETC:          e.ParentLogger.ETC.BranchSub(),
			Msgtags:      e.subLoggerInitChain.msgtags,
		}
		return logger
	} else {
		if e.ParentLogger.RootContext.PanicOnMisuse {
			panic("5a2d0286d0bd4adca2ceda0f68f3b0fb")
		} else {
			return nil
		}
	}
}
