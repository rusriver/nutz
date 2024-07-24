package logger

import "github.com/rs/zerolog"

var _ ILogger = (*Logger)(nil)

type Settings struct {
	PanicOnMisuse bool
	MsgtagKey     string
	OnSendHook    func(e *Event) (doSend bool)
}

type Logger struct {
	Settings     *Settings
	TheSubsystem string // per logger

	zeroLogger *zerolog.Logger
}

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

func (l *Logger) Debug() (e IEvent) {
	return &Event{
		TheLevel:     zerolog.DebugLevel.String(),
		ParentLogger: l,
		zerologEvent: l.zeroLogger.Debug(),
	}
}

func (l *Logger) Trace() (e IEvent) {
	return &Event{
		TheLevel:     zerolog.TraceLevel.String(),
		ParentLogger: l,
		zerologEvent: l.zeroLogger.Trace(),
	}
}

func (l *Logger) Error() (e IEvent) {
	return &Event{
		TheLevel:     zerolog.ErrorLevel.String(),
		ParentLogger: l,
		zerologEvent: l.zeroLogger.Error(),
	}
}

func (l *Logger) Err(err error) (e IEvent) {
	return &Event{
		TheLevel:     zerolog.ErrorLevel.String(),
		ParentLogger: l,
		zerologEvent: l.zeroLogger.Err(err),
	}
}

func (l *Logger) Warn() (e IEvent) {
	return &Event{
		TheLevel:     zerolog.WarnLevel.String(),
		ParentLogger: l,
		zerologEvent: l.zeroLogger.Warn(),
	}
}

func (l *Logger) Info() (e IEvent) {
	return &Event{
		TheLevel:     zerolog.InfoLevel.String(),
		ParentLogger: l,
		zerologEvent: l.zeroLogger.Info(),
	}
}

// Can be used if you intend to call ILogger() on it later
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
