package logger

import (
	"fmt"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rusriver/nutz/logger/loggerapi"
)

type subLoggerInitChain struct {
	zerologContext *zerolog.Context
	isInactive     bool // for new sub-logger
	traceContext   *loggerapi.TraceContext
	rctx           *loggerapi.RootContext
	uo             any
	accumulator    *Accumulator
	interceptor    *Interceptor
}

var _ loggerapi.ISubLoggerInitChain = (*subLoggerInitChain)(nil)

func (slich *subLoggerInitChain) Inactive() loggerapi.ISubLoggerInitChain {
	slich.isInactive = true
	return slich
}

func (slich *subLoggerInitChain) Caller(skip ...int) loggerapi.ISubLoggerInitChain {
	skipCount := 4
	if len(skip) > 0 {
		skipCount = skip[0]
	}
	zc := slich.zerologContext.CallerWithSkipFrameCount(skipCount)
	slich.zerologContext = &zc
	return slich
}

func (slich *subLoggerInitChain) Str(k, v string) loggerapi.ISubLoggerInitChain {
	zc := slich.zerologContext.Str(k, v)
	slich.zerologContext = &zc
	return slich
}

func (slich *subLoggerInitChain) Strf(k, s string, vv ...any) loggerapi.ISubLoggerInitChain {
	zc := slich.zerologContext.Str(k, fmt.Sprintf(s, vv...))
	slich.zerologContext = &zc
	return slich
}

func (slich *subLoggerInitChain) Strs(k string, vv []string) loggerapi.ISubLoggerInitChain {
	zc := slich.zerologContext.Strs(k, vv)
	slich.zerologContext = &zc
	return slich
}

func (slich *subLoggerInitChain) Time(k string, t time.Time) loggerapi.ISubLoggerInitChain {
	zc := slich.zerologContext.Time(k, t)
	slich.zerologContext = &zc
	return slich
}

func (slich *subLoggerInitChain) Int(k string, v int) loggerapi.ISubLoggerInitChain {
	zc := slich.zerologContext.Int(k, v)
	slich.zerologContext = &zc
	return slich
}

func (slich *subLoggerInitChain) Ints(k string, v []int) loggerapi.ISubLoggerInitChain {
	zc := slich.zerologContext.Ints(k, v)
	slich.zerologContext = &zc
	return slich
}

func (slich *subLoggerInitChain) Float64(k string, v float64) loggerapi.ISubLoggerInitChain {
	zc := slich.zerologContext.Float64(k, v)
	slich.zerologContext = &zc
	return slich
}

func (slich *subLoggerInitChain) Floats64(k string, v []float64) loggerapi.ISubLoggerInitChain {
	zc := slich.zerologContext.Floats64(k, v)
	slich.zerologContext = &zc
	return slich
}

func (slich *subLoggerInitChain) Array(k string, v *zerolog.Array) loggerapi.ISubLoggerInitChain {
	zc := slich.zerologContext.Array(k, v)
	slich.zerologContext = &zc
	return slich
}

func (slich *subLoggerInitChain) Dict(k string, v *zerolog.Event) loggerapi.ISubLoggerInitChain {
	zc := slich.zerologContext.Dict(k, v)
	slich.zerologContext = &zc
	return slich
}

func (slich *subLoggerInitChain) RawJSON(k string, rj []byte) loggerapi.ISubLoggerInitChain {
	zc := slich.zerologContext.RawJSON(k, rj)
	slich.zerologContext = &zc
	return slich
}

func (slich *subLoggerInitChain) Bytes(k string, bb []byte) loggerapi.ISubLoggerInitChain {
	zc := slich.zerologContext.Bytes(k, bb)
	slich.zerologContext = &zc
	return slich
}

func (slich *subLoggerInitChain) WithClonedRootContext(optFuncs ...func(rctx *loggerapi.RootContext)) loggerapi.ISubLoggerInitChain {
	rctx := slich.rctx.Clone()
	for _, f := range optFuncs {
		f(rctx)
	}
	if rctx.MainWriter == nil {
		rctx.MainWriter = os.Stderr
	}

	slich.interceptor = &Interceptor{rctx: rctx}

	newCtx := slich.zerologContext.Logger().Output(slich.interceptor).With()
	slich.zerologContext = &newCtx

	slich.rctx = rctx

	return slich
}

func (slich *subLoggerInitChain) ILogger() loggerapi.ILogger {
	zeroLogger := slich.zerologContext.Logger()
	logger := &Logger{
		RootContext: slich.rctx,
		IsInactive:  slich.isInactive, // inherit from the init chain
		UserObject:  slich.uo,
		ZeroLogger:  &zeroLogger,
		TraceContext: &loggerapi.TraceContext{
			ETC:     slich.traceContext.ETC,
			Msgtags: slich.traceContext.Msgtags,
		},
		accumulator: slich.accumulator,
		interceptor: slich.interceptor,
	}
	return logger
}

func (slich *subLoggerInitChain) AddUserObject(uo any) loggerapi.ISubLoggerInitChain {
	slich.uo = uo
	return slich
}
