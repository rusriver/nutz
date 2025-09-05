package controlflow

import "github.com/rusriver/nutz/loggeterr-api/trace"

/*
This works, but lacks the ID information passed.

func Must[T any](result T, err error) T {
	if err != nil {
		RaiseCause(nil, &trace.BacktracePoint{"123", err.Error()})
	}
	return result
}
*/

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

type MustPair[T any] struct {
	result T
	err    error
}

// must-value-error-pair
func Mustvep[T any](result T, err error) MustPair[T] {
	return MustPair[T]{result, err}
}

/*
This actually works:

	var ft *trace.Forwardtrace

	values := controlflow.Must(controlflow.Mustvep(
		url.ParseQuery("a=1&s=2;"),
	), ft, "b17b648bbee12b797619")

This is the alternative to do the same explicitly:

	values, err := url.ParseQuery("a=1&s=2;")
	if err != nil {
		controlflow.RaiseCause(ft, &trace.BacktracePoint{"b17b648bbee12b797619", err.Error()})
	}

Good thing about first variant, is that it quite easily can be wrapped into a syntax sugar:

	values := Must(url.ParseQuery("a=1&s=2;"), ft, "b17b648bbee12b797619")
*/
func Must[T any](mp MustPair[T], ft *trace.Forwardtrace, id string) T {
	if mp.err != nil {
		RaiseCause(ft, &trace.BacktracePoint{id, mp.err.Error()})
	}
	return mp.result
}
