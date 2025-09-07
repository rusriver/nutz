package dslre

import (
	"os"

	"github.com/rusriver/nutz/controlflow"
	loggeterrapi "github.com/rusriver/nutz/loggeterr-api"
	"github.com/rusriver/nutz/loggeterr-api/trace"
)

type Registry[T_parsed_program any, T_arg any, T_result any] struct {
	InternalTracer                 loggeterrapi.ITracer
	InternalCache                  ICache[T_parsed_program]
	InternalParserFunc             func(literal string) (pp T_parsed_program)
	InternalPackerFunc0            func(pp T_parsed_program) (res T_result)
	InternalPackerFunc1            func(pp T_parsed_program, arg T_arg) (res T_result)
	InternalPackerFuncN            func(pp T_parsed_program, args ...T_arg) (res T_result)
	InternalIderFunc               func(literal string) (id string)
	InternalThresholdLiteralLength int
	InternalStorage                IReadFile
}

func NewRegistry[T_parsed_program any, T_arg any, T_result any](
	tracer loggeterrapi.ITracer,
	of func(*Registry[T_parsed_program, T_arg, T_result]),
) *Registry[T_parsed_program, T_arg, T_result] {

	r := &Registry[T_parsed_program, T_arg, T_result]{}
	// defaults here
	r.InternalTracer = tracer
	r.InternalThresholdLiteralLength = 60

	// user-side
	of(r)

	// validation
	if r.InternalCache == nil {
		controlflow.RaiseCause(tracer.GetForwardtrace(), &trace.BacktracePoint{"7aa9b11a0ca7bc", "cache not set"})
	}
	if r.InternalParserFunc == nil {
		controlflow.RaiseCause(tracer.GetForwardtrace(), &trace.BacktracePoint{"f45c318d6609", "InternalParserFunc not set"})
	}

	return r
}

func (r *Registry[T_parsed_program, T_arg, T_result]) getParsedProgramWithAutoIdFromLiteral(literal string) (pp T_parsed_program) {
	var id string
	if len(literal) <= r.InternalThresholdLiteralLength {
		id = literal
	} else {
		if r.InternalIderFunc == nil {
			controlflow.RaiseCause(r.InternalTracer.GetForwardtrace(),
				trace.BacktracePointf("9df89e04bdfa", "InternalIderFunc not set, while the literal is longer than %v character",
					r.InternalThresholdLiteralLength))
		}
		id = r.InternalIderFunc(literal)
	}
	pp, ok := r.InternalCache.Get(id)
	if !ok {
		pp = r.InternalParserFunc(literal)
		r.InternalCache.Set(id, pp)
	}
	return pp
}

func (r *Registry[T_parsed_program, T_arg, T_result]) getParsedProgramWithExplicitIdFromLiteral(id string, literal string) (pp T_parsed_program) {
	pp, ok := r.InternalCache.Get(id)
	if !ok {
		pp = r.InternalParserFunc(literal)
		r.InternalCache.Set(id, pp)
	}
	return pp
}

func (r *Registry[T_parsed_program, T_arg, T_result]) getParsedProgramWithAutoIdFromStorage(path string) (pp T_parsed_program) {
	id := path
	pp, ok := r.InternalCache.Get(id)
	if !ok {
		literal := r.getLiteralFromStorage(path)
		pp = r.InternalParserFunc(literal)
		r.InternalCache.Set(id, pp)
	}
	return pp
}

func (r *Registry[T_parsed_program, T_arg, T_result]) getParsedProgramWithExplicitIdFromStorage(id string, path string) (pp T_parsed_program) {
	pp, ok := r.InternalCache.Get(id)
	if !ok {
		literal := r.getLiteralFromStorage(path)
		pp = r.InternalParserFunc(literal)
		r.InternalCache.Set(id, pp)
	}
	return pp
}

func (r *Registry[T_parsed_program, T_arg, T_result]) getLiteralFromStorage(path string) (literal string) {
	storage := r.InternalStorage
	if storage == nil {
		literal = string(
			controlflow.VEP(os.ReadFile(path)).MustWith(r.InternalTracer.GetForwardtrace(), "bd60d0022fe4"),
		)
	} else {
		literal = string(
			controlflow.VEP(storage.ReadFile(path)).MustWith(r.InternalTracer.GetForwardtrace(), "508fa97a3490"),
		)
	}
	return
}
