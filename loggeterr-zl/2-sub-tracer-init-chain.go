package loggeterrzl

import (
	"fmt"
	"os"
	"time"

	"github.com/rs/zerolog"
	loggeterrapi "github.com/rusriver/nutz/loggeterr-api"
	"github.com/rusriver/nutz/loggeterr-api/trace"
)

type subTracerInitChain struct {
	zerologContext *zerolog.Context
	isInactive     bool // for new sub-logger
	forwardtrace   *trace.Forwardtrace
	rctx           *loggeterrapi.RootContext
	uo             any
	accumulator    *Accumulator
	interceptor    *Interceptor
}

var _ loggeterrapi.ISubTracerInitChain = (*subTracerInitChain)(nil)

func (stic *subTracerInitChain) Inactive() loggeterrapi.ISubTracerInitChain {
	stic.isInactive = true
	return stic
}

func (stic *subTracerInitChain) Caller(skip ...int) loggeterrapi.ISubTracerInitChain {
	skipCount := 4
	if len(skip) > 0 {
		skipCount = skip[0]
	}
	zc := stic.zerologContext.CallerWithSkipFrameCount(skipCount)
	stic.zerologContext = &zc
	return stic
}

func (stic *subTracerInitChain) Str(k, v string) loggeterrapi.ISubTracerInitChain {
	zc := stic.zerologContext.Str(k, v)
	stic.zerologContext = &zc
	return stic
}

func (stic *subTracerInitChain) Strf(k, s string, vv ...any) loggeterrapi.ISubTracerInitChain {
	zc := stic.zerologContext.Str(k, fmt.Sprintf(s, vv...))
	stic.zerologContext = &zc
	return stic
}

func (stic *subTracerInitChain) Strs(k string, vv []string) loggeterrapi.ISubTracerInitChain {
	zc := stic.zerologContext.Strs(k, vv)
	stic.zerologContext = &zc
	return stic
}

func (stic *subTracerInitChain) Time(k string, t time.Time) loggeterrapi.ISubTracerInitChain {
	zc := stic.zerologContext.Time(k, t)
	stic.zerologContext = &zc
	return stic
}

func (stic *subTracerInitChain) Int(k string, v int) loggeterrapi.ISubTracerInitChain {
	zc := stic.zerologContext.Int(k, v)
	stic.zerologContext = &zc
	return stic
}

func (stic *subTracerInitChain) Ints(k string, v []int) loggeterrapi.ISubTracerInitChain {
	zc := stic.zerologContext.Ints(k, v)
	stic.zerologContext = &zc
	return stic
}

func (stic *subTracerInitChain) Float64(k string, v float64) loggeterrapi.ISubTracerInitChain {
	zc := stic.zerologContext.Float64(k, v)
	stic.zerologContext = &zc
	return stic
}

func (stic *subTracerInitChain) Floats64(k string, v []float64) loggeterrapi.ISubTracerInitChain {
	zc := stic.zerologContext.Floats64(k, v)
	stic.zerologContext = &zc
	return stic
}

func (stic *subTracerInitChain) Array(k string, v *zerolog.Array) loggeterrapi.ISubTracerInitChain {
	zc := stic.zerologContext.Array(k, v)
	stic.zerologContext = &zc
	return stic
}

func (stic *subTracerInitChain) Dict(k string, v *zerolog.Event) loggeterrapi.ISubTracerInitChain {
	zc := stic.zerologContext.Dict(k, v)
	stic.zerologContext = &zc
	return stic
}

func (stic *subTracerInitChain) RawStruct(k string, rj []byte) loggeterrapi.ISubTracerInitChain {
	zc := stic.zerologContext.RawJSON(k, rj)
	stic.zerologContext = &zc
	return stic
}

func (stic *subTracerInitChain) Bytes(k string, bb []byte) loggeterrapi.ISubTracerInitChain {
	zc := stic.zerologContext.Bytes(k, bb)
	stic.zerologContext = &zc
	return stic
}

func (stic *subTracerInitChain) WithClonedRootContext(optFuncs ...func(rctx *loggeterrapi.RootContext)) loggeterrapi.ISubTracerInitChain {
	rctx := stic.rctx.Clone()
	for _, f := range optFuncs {
		f(rctx)
	}
	if rctx.MainWriter == nil {
		rctx.MainWriter = os.Stderr
	}

	stic.interceptor = &Interceptor{rctx: rctx}

	newCtx := stic.zerologContext.Logger().Output(stic.interceptor).With()
	stic.zerologContext = &newCtx

	stic.rctx = rctx

	return stic
}

func (stic *subTracerInitChain) ITracer() loggeterrapi.ITracer {
	zeroLogger := stic.zerologContext.Logger()
	tracer := &Tracer{
		RootContext:   stic.rctx,
		IsInactive:    stic.isInactive, // inherit from the init chain
		UserObject:    stic.uo,
		ZeroLoggerPtr: &zeroLogger,
		Forwardtrace: &trace.Forwardtrace{
			ETC:  stic.forwardtrace.ETC,
			Tags: stic.forwardtrace.Tags,
		},
		accumulator: stic.accumulator,
		interceptor: stic.interceptor,
	}
	return tracer
}

func (stic *subTracerInitChain) SetUserObject(uo any) loggeterrapi.ISubTracerInitChain {
	stic.uo = uo
	return stic
}

func (stic *subTracerInitChain) SetForwardtrace(ft *trace.Forwardtrace) loggeterrapi.ISubTracerInitChain {
	stic.forwardtrace = ft
	return stic
}
