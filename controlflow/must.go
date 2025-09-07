package controlflow

import (
	"runtime"
	"strconv"

	"github.com/rusriver/nutz/loggeterr-api/trace"
)

// This works, but lacks FT and ID, instead uses source code file and LN
// to tag error. Good for simple and local things, for serious production
// use see VEP tools.
func Must[T any](result T, err error) T {
	if err != nil {
		var id string

		{
			maxLen := 12
			_, file, line, ok := runtime.Caller(3)
			if ok {
				if len(file) >= maxLen {
					file = "..." + file[len(file)-maxLen:]
				}
				id += file + ":" + strconv.Itoa(line)
			}
		}
		if len(id) > 0 {
			id += ".."
		}
		{
			maxLen := 22
			_, file, line, ok := runtime.Caller(2)
			if ok {
				if len(file) >= maxLen {
					file = "..." + file[len(file)-maxLen:]
				}
				id += file + ":" + strconv.Itoa(line)
			}
		}
		if len(id) > 0 {
			id += ".."
		}
		{
			maxLen := 33
			_, file, line, ok := runtime.Caller(1)
			if ok {
				if len(file) >= maxLen {
					file = "..." + file[len(file)-maxLen:]
				}
				id += file + ":" + strconv.Itoa(line)
			}
		}

		RaiseCause(nil, &trace.BacktracePoint{id, err.Error()})
	}
	return result
}

/*
type MustObject struct {
	ft *trace.Forwardtrace
	id string
}

func MakeMustObject(ft *trace.Forwardtrace, id string) MustObject {
	return MustObject{ft, id}
}

This doesn't work in Go, because "method must have no type parameters"...

func (mo MustObject) M2[T any](result T, err error) T {
	if err != nil {
		RaiseCause(ft, &trace.BacktracePoint{id, err.Error()})
	}
	return result
}
*/

/*
This doesn't work in Go either, because you can't write:
	values := controlflow.Must(url.ParseQuery("a=1&s=2;"), nil, "123123")

func Must[T any](result T, err error, ft *trace.Forwardtrace, id string) T {
	if err != nil {
		RaiseCause(nil, &trace.BacktracePoint{"123", err.Error()})
	}
	return result
}
*/

// Value-Error pair
type VEP_Type[T any] struct {
	result T
	err    error
}

// Execute the subject function, and make a value-error pair object
func VEP[T any](result T, err error) VEP_Type[T] {
	return VEP_Type[T]{result, err}
}

/*
This actually works:

	var ft *trace.Forwardtrace

	values := controlflow.MustVep(controlflow.VEP(
		url.ParseQuery("a=1&s=2;"),
	), ft, "b17b648bbee12b797619")

This is the alternative to do the same explicitly:

	values, err := url.ParseQuery("a=1&s=2;")
	if err != nil {
		controlflow.RaiseCause(ft, &trace.BacktracePoint{"b17b648bbee12b797619", err.Error()})
	}

Good thing about first variant, is that it quite easily can be wrapped into a syntax sugar:

	values := MustVep(url.ParseQuery("a=1&s=2;"), ft, "b17b648bbee12b797619")
*/
func MustVep[T any](vep VEP_Type[T], ft *trace.Forwardtrace, id string) T {
	if vep.err != nil {
		RaiseCause(ft, &trace.BacktracePoint{id, vep.err.Error()})
	}
	return vep.result
}

// -----------

func (vep VEP_Type[T]) MustWith(ft *trace.Forwardtrace, id string) T {
	if vep.err != nil {
		RaiseCause(ft, &trace.BacktracePoint{id, vep.err.Error()})
	}
	return vep.result
}
