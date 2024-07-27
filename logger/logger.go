package logger

import "github.com/rs/zerolog"

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
		zeroLogger: zl,
	}
	return logger
}

func (l *Logger) DebugEvent() (e IEvent) {
	return &Event{
		IsInactive:   l.IsInactive,
		TheLevel:     zerolog.DebugLevel.String(),
		ParentLogger: l,
		zerologEvent: l.zeroLogger.Debug(),
	}
}

func (l *Logger) TraceEvent() (e IEvent) {
	return &Event{
		IsInactive:   l.IsInactive,
		TheLevel:     zerolog.TraceLevel.String(),
		ParentLogger: l,
		zerologEvent: l.zeroLogger.Trace(),
	}
}

func (l *Logger) ErrorEvent() (e IEvent) {
	return &Event{
		IsInactive:   l.IsInactive,
		TheLevel:     zerolog.ErrorLevel.String(),
		ParentLogger: l,
		zerologEvent: l.zeroLogger.Error(),
	}
}

func (l *Logger) ErrEvent(err error) (e IEvent) {
	return &Event{
		IsInactive:   l.IsInactive,
		TheLevel:     zerolog.ErrorLevel.String(),
		ParentLogger: l,
		zerologEvent: l.zeroLogger.Err(err),
	}
}

func (l *Logger) WarnEvent() (e IEvent) {
	return &Event{
		IsInactive:   l.IsInactive,
		TheLevel:     zerolog.WarnLevel.String(),
		ParentLogger: l,
		zerologEvent: l.zeroLogger.Warn(),
	}
}

func (l *Logger) InfoEvent() (e IEvent) {
	return &Event{
		IsInactive:   l.IsInactive,
		TheLevel:     zerolog.InfoLevel.String(),
		ParentLogger: l,
		zerologEvent: l.zeroLogger.Info(),
	}
}

// Can be used if you intend to call ILogger() on it later.
// By default, all new sub-loggers are active, even if parent was inactivated.
func (l *Logger) SubLoggerInitChain() IEvent {
	zerologContext := l.zeroLogger.With()
	e := &Event{
		subLoggerInitChain: &subLoggerInitChain{
			zerologContext: &zerologContext,
		},
		ParentLogger: l,
	}
	return e
}

func (l *Logger) GetZeroLoggerPtr() *zerolog.Logger {
	return l.zeroLogger
}

// This is useful to access Settings later after init.
func (l *Logger) GetNutzLogger() *Logger {
	return l
}

func (l *Logger) SetLevel(level zerolog.Level) ILogger {
	return l
}

// Attach any user object to the event object, to be used later
// inside the callback hooks. Intended to pass user-level runtime data to the
// hook(s). Doesn't add anything to the logline serializations.
// Normally, the UO must be a pointer, and at the receiving side in the hook
// you do something along the lines:
//
//	myObject, _ := e.ParentLogger.UserObject.(*MyObject)
//	if myObject !=  nil {
//		...
//	}
func (l *Logger) AddUO(uo any) ILogger {
	l.UserObject = uo
	return l
}
