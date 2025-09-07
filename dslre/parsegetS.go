package dslre

import (
	"github.com/rusriver/nutz/controlflow"
	"github.com/rusriver/nutz/loggeterr-api/trace"
)

// Program from storage, automatic Id, no args
func (r *Registry[T_parsed_program, T_arg, T_result]) ParsegetS(path string) (res T_result) {
	packerFunc := r.InternalPackerFunc0
	if packerFunc == nil {
		controlflow.RaiseCause(r.InternalTracer.GetForwardtrace(), &trace.BacktracePoint{"012737f45d07", "InternalPackerFunc0 not set"})
	}
	pp := r.getParsedProgramWithAutoIdFromStorage(path)
	res = packerFunc(pp)
	return
}

// Program from storage, automatic Id, one arg
func (r *Registry[T_parsed_program, T_arg, T_result]) ParsegetS1(path string, arg1 T_arg) (res T_result) {
	packerFunc := r.InternalPackerFunc1
	if packerFunc == nil {
		controlflow.RaiseCause(r.InternalTracer.GetForwardtrace(), &trace.BacktracePoint{"67300d800635", "InternalPackerFunc1 not set"})
	}
	pp := r.getParsedProgramWithAutoIdFromStorage(path)
	res = packerFunc(pp, arg1)
	return
}

// Program from storage, automatic Id, variable number of args
func (r *Registry[T_parsed_program, T_arg, T_result]) ParsegetSN(path string, args ...T_arg) (res T_result) {
	packerFunc := r.InternalPackerFuncN
	if packerFunc == nil {
		controlflow.RaiseCause(r.InternalTracer.GetForwardtrace(), &trace.BacktracePoint{"580ddac2ccab", "InternalPackerFuncN not set"})
	}
	pp := r.getParsedProgramWithAutoIdFromStorage(path)
	res = packerFunc(pp, args...)
	return
}

// Program from storage, explicit Id, no args
func (r *Registry[T_parsed_program, T_arg, T_result]) ParsegetSi(id string, path string) (res T_result) {
	packerFunc := r.InternalPackerFunc0
	if packerFunc == nil {
		controlflow.RaiseCause(r.InternalTracer.GetForwardtrace(), &trace.BacktracePoint{"e63ce8803fee", "InternalPackerFunc0 not set"})
	}
	pp := r.getParsedProgramWithExplicitIdFromStorage(id, path)
	res = packerFunc(pp)
	return
}

// Program from storage, explicit Id, one arg
func (r *Registry[T_parsed_program, T_arg, T_result]) ParsegetS1i(id string, path string, arg1 T_arg) (res T_result) {
	packerFunc := r.InternalPackerFunc1
	if packerFunc == nil {
		controlflow.RaiseCause(r.InternalTracer.GetForwardtrace(), &trace.BacktracePoint{"a711c2803fe2", "InternalPackerFunc1 not set"})
	}
	pp := r.getParsedProgramWithExplicitIdFromStorage(id, path)
	res = packerFunc(pp, arg1)
	return
}

// Program from storage, explicit Id, variable number of args
func (r *Registry[T_parsed_program, T_arg, T_result]) ParsegetSNi(id string, path string, args ...T_arg) (res T_result) {
	packerFunc := r.InternalPackerFuncN
	if packerFunc == nil {
		controlflow.RaiseCause(r.InternalTracer.GetForwardtrace(), &trace.BacktracePoint{"30b110cfcec0", "InternalPackerFuncN not set"})
	}
	pp := r.getParsedProgramWithExplicitIdFromStorage(id, path)
	res = packerFunc(pp, args...)
	return
}
