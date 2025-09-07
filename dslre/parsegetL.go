package dslre

import (
	"github.com/rusriver/nutz/controlflow"
	"github.com/rusriver/nutz/loggeterr-api/trace"
)

// Literal program, automatic Id, no args
func (r *Registry[T_parsed_program, T_arg, T_result]) ParsegetL(literal string) (res T_result) {
	packerFunc := r.InternalPackerFunc0
	if packerFunc == nil {
		controlflow.RaiseCause(r.InternalTracer.GetForwardtrace(), &trace.BacktracePoint{"abe30c69186f", "InternalPackerFunc0 not set"})
	}
	pp := r.getParsedProgramWithAutoIdFromLiteral(literal)
	res = packerFunc(pp)
	return
}

// Literal program, automatic Id, one arg
func (r *Registry[T_parsed_program, T_arg, T_result]) ParsegetL1(literal string, arg1 T_arg) (res T_result) {
	packerFunc := r.InternalPackerFunc1
	if packerFunc == nil {
		controlflow.RaiseCause(r.InternalTracer.GetForwardtrace(), &trace.BacktracePoint{"ed6298da237e", "InternalPackerFunc1 not set"})
	}
	pp := r.getParsedProgramWithAutoIdFromLiteral(literal)
	res = packerFunc(pp, arg1)
	return
}

// Literal program, automatic Id, variable number of args
func (r *Registry[T_parsed_program, T_arg, T_result]) ParsegetLN(literal string, args ...T_arg) (res T_result) {
	packerFunc := r.InternalPackerFuncN
	if packerFunc == nil {
		controlflow.RaiseCause(r.InternalTracer.GetForwardtrace(), &trace.BacktracePoint{"d7c638a5ba5f", "InternalPackerFuncN not set"})
	}
	pp := r.getParsedProgramWithAutoIdFromLiteral(literal)
	res = packerFunc(pp, args...)
	return
}

// Literal program, explicit Id, no args
func (r *Registry[T_parsed_program, T_arg, T_result]) ParsegetLi(id string, literal string) (res T_result) {
	packerFunc := r.InternalPackerFunc0
	if packerFunc == nil {
		controlflow.RaiseCause(r.InternalTracer.GetForwardtrace(), &trace.BacktracePoint{"7845d2d511a7", "InternalPackerFunc0 not set"})
	}
	pp := r.getParsedProgramWithExplicitIdFromLiteral(id, literal)
	res = packerFunc(pp)
	return
}

// Literal program, explicit Id, one arg
func (r *Registry[T_parsed_program, T_arg, T_result]) ParsegetL1i(id string, literal string, arg1 T_arg) (res T_result) {
	packerFunc := r.InternalPackerFunc1
	if packerFunc == nil {
		controlflow.RaiseCause(r.InternalTracer.GetForwardtrace(), &trace.BacktracePoint{"73b54a1bc843", "InternalPackerFunc1 not set"})
	}
	pp := r.getParsedProgramWithExplicitIdFromLiteral(id, literal)
	res = packerFunc(pp, arg1)
	return
}

// Literal program, explicit Id, variable number of args
func (r *Registry[T_parsed_program, T_arg, T_result]) ParsegetLNi(id string, literal string, args ...T_arg) (res T_result) {
	packerFunc := r.InternalPackerFuncN
	if packerFunc == nil {
		controlflow.RaiseCause(r.InternalTracer.GetForwardtrace(), &trace.BacktracePoint{"9228b56075fd", "InternalPackerFuncN not set"})
	}
	pp := r.getParsedProgramWithExplicitIdFromLiteral(id, literal)
	res = packerFunc(pp, args...)
	return
}
