package logger

import (
	"errors"

	"github.com/rs/zerolog"
)

// An initial zerologger must be provided by yours explicitly.
// Typical idiom:
/*
	import(
		"github.com/rs/zerolog/log"
	)
		...
		my.Logger = logger.New(&log.Logger, func(s *logger.Settings) {
			s.PanicOnMisuse = false
		})
*/
func New(zl *zerolog.Logger, optFuncs ...func(s *Settings)) (logger *Logger) {
	settings := &Settings{
		PanicOnMisuse: true,
		MsgtagKey:     "msgtags",
	}
	for _, f := range optFuncs {
		f(settings)
	}
	logger = &Logger{
		Settings:   settings,
		ZeroLogger: zl,
	}
	return logger
}

func (l *Logger) DebugEvent() (e IEvent) {
	return &Event{
		IsInactive:   l.IsInactive,
		TheLevel:     zerolog.DebugLevel.String(),
		ParentLogger: l,
		zerologEvent: l.ZeroLogger.Debug(),
	}
}

func (l *Logger) TraceEvent() (e IEvent) {
	return &Event{
		IsInactive:   l.IsInactive,
		TheLevel:     zerolog.TraceLevel.String(),
		ParentLogger: l,
		zerologEvent: l.ZeroLogger.Trace(),
	}
}

func (l *Logger) ErrorEvent() (e IEvent) {
	return &Event{
		IsInactive:        l.IsInactive,
		TheLevel:          zerolog.ErrorLevel.String(),
		ParentLogger:      l,
		zerologEvent:      l.ZeroLogger.Error(),
		TheTitle:          "ERROR", // the default title
		theTitleIsDefault: true,
	}
}

func (l *Logger) ErrEvent(err error) (e IEvent) {
	if err == nil {
		err = errors.New("<nil>")
	}
	return &Event{
		IsInactive:        l.IsInactive,
		TheLevel:          zerolog.ErrorLevel.String(),
		ParentLogger:      l,
		zerologEvent:      l.ZeroLogger.Err(err),
		TheTitle:          "ERROR", // the default title
		theTitleIsDefault: true,
	}
}

func (l *Logger) WarnEvent() (e IEvent) {
	return &Event{
		IsInactive:        l.IsInactive,
		TheLevel:          zerolog.WarnLevel.String(),
		ParentLogger:      l,
		zerologEvent:      l.ZeroLogger.Warn(),
		TheTitle:          "WARNING", // the default title
		theTitleIsDefault: true,
	}
}

func (l *Logger) InfoEvent() (e IEvent) {
	return &Event{
		IsInactive:   l.IsInactive,
		TheLevel:     zerolog.InfoLevel.String(),
		ParentLogger: l,
		zerologEvent: l.ZeroLogger.Info(),
	}
}

func (l *Logger) SubLoggerInitChain() IEvent {
	zerologContext := l.ZeroLogger.With()
	e := &Event{
		subLoggerInitChain: &subLoggerInitChain{
			zerologContext: &zerologContext,
		},
		ParentLogger: l,
	}
	return e
}

func (l *Logger) GetZeroLoggerPtr() *zerolog.Logger {
	return l.ZeroLogger
}

func (l *Logger) GetNutzLogger() *Logger {
	return l
}

func (l *Logger) SetLevel(level zerolog.Level) ILogger {
	return l
}

func (l *Logger) AddUO(uo any) ILogger {
	l.UserObject = uo
	return l
}
