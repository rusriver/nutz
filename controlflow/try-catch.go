package controlflow

import (
	"fmt"
	"strings"

	"github.com/rusriver/nutz/loggeterr/loggeterrapi/trace"
)

type Event struct {
	RetryAttemptNb int
	Forwardtrace   *trace.Forwardtrace
	Backtrace      FullBacktrace
}

// First element is most specific, last ones are least specific - this is the inverse of forwardtrace tags;
type FullBacktrace []*BacktraceDomainChain

// First element is most specific, last ones are least specific - this is the inverse of forwardtrace tags;
type BacktraceDomainChain []*trace.BacktracePoint

/*
	Forwardtrace has last elements most specific; backtrace is the inverse, the most
	specific elements are first ones.

	To get most specific elements of forwardtrace, use Last() method on Tags field.

	To get most specific elements of the backtrace, just get [0] element, off the domains
	array, and off the BPs array.

	Examples:

	Get the cause (most specific) BP, in the current domain:
		bp = event.Backtrace.Last(-1)[0]

	Get the very root cause:
		bp = event.Backtrace.Last[0][0]
*/

// E.g. to get the one before last, use Last(-2);
// The purpose is to easily get the last thing. Also used internally
func (b FullBacktrace) Last(neg_i int) *BacktraceDomainChain {
	if neg_i >= 0 {
		return &BacktraceDomainChain{}
	}
	i := len(b) + neg_i
	if i > 0 {
		return b[i]
	} else {
		return &BacktraceDomainChain{}
	}
}

// E.g. to get the one before last, use Last(-2);
// The purpose is to easily get the last thing. Also used internally
func (b BacktraceDomainChain) Last(neg_i int) *trace.BacktracePoint {
	if neg_i >= 0 {
		return nil
	}
	i := len(b) + neg_i
	if i > 0 {
		return b[i]
	} else {
		return nil
	}
}

var _ error = (*Event)(nil)

func (e *Event) Error() string {
	return e.String()
}

var _ fmt.Stringer = (*Event)(nil)

func (e *Event) String() string {
	/*
		Example string:
			{: Forwardtrace {: ETC [1 2 3 4 5], Tags [qq ww ee]}, Backtrace {[123123 "detail"] [123123 "detail"] [123123 "detail"], [123123 "detail"] [123123 "detail"]}}
	*/
	br := strings.Builder{}
	for i := len(e.Backtrace) - 1; i >= 0; i-- {
		// x := e.Backtrace[i]
		// switch xv := x.(type) {
		// case error:
		// 	br.WriteString(xv.Error())
		// 	br.WriteString("; ")
		// case string:
		// 	br.WriteString(xv)
		// 	br.WriteString("; ")
		// default:
		// 	br.WriteString(fmt.Sprintf("%v; ", xv))
		// }
	}
	// br.WriteString(
	// 	e.TheCauseError.Error(),
	// )
	return br.String()
}

func Try(f func() (err error)) (e *Event) {
	defer func() {
		// if x := recover(); x != nil {
		// 	switch xv := x.(type) {
		// 	case *Event:
		// 		e = xv
		// 	case error:
		// 		e = &Event{
		// 			Cause:         x,
		// 			TheCauseError: xv,
		// 		}
		// 	default:
		// 		e = &Event{
		// 			Cause:         x,
		// 			TheCauseError: fmt.Errorf("%v", xv),
		// 		}
		// 	}
		// }
	}()
	err := f()
	if err != nil {
		// e = &Event{
		// 	Cause:         err,
		// 	TheCauseError: err,
		// }
	}
	return
}

func (e *Event) Catch(f func(e *Event)) *Event {
	if e != nil {
		f(e)
	}
	return e
}

func (e *Event) OK(f func(e *Event)) *Event {
	if e == nil {
		f(e)
	}
	return e
}

// Allows to try and retry an operation several times.
func Retry(f func() (err error), f2 func(e *Event) (end bool)) (e *Event) {
	attemptNb := 0
	for {
		attemptNb++
		e = Try(func() (err error) {
			err = f()
			return
		})
		if e == nil {
			break
		} else {
			e.RetryAttemptNb = attemptNb
			end := f2(e)
			if end {
				break
			}
		}
	}
	return
}

func (existing_event *Event) Reraise(added_backtrace_point *trace.BacktracePoint) {
	existing_event.Wrap(added_backtrace_point)
	panic(existing_event)
}

// Example usage: controlflow.RaiseCause(tracer.GetForwardtrace(), trace.BacktracePointf("123123", "the cause"))
// controlflow.RaiseCause(nil, trace.BacktracePointf("123123", "the cause"))
func RaiseCause(forward_trace *trace.Forwardtrace, first_backtrace_point *trace.BacktracePoint) {
	event := NewEvent(forward_trace, first_backtrace_point)
	panic(event)
}

// Example usage: controlflow.RaiseCause(tracer.GetForwardtrace(), trace.BacktracePointf("123123", "the cause"))
// controlflow.RaiseCause(nil, trace.BacktracePointf("123123", "the cause"))
func NewEvent(forward_trace *trace.Forwardtrace, first_backtrace_point *trace.BacktracePoint) *Event {
	if forward_trace == nil {
		forward_trace = &trace.Forwardtrace{}
	}
	event := &Event{
		Forwardtrace: forward_trace,
		Backtrace: FullBacktrace{
			&BacktraceDomainChain{},
		},
	}
	event.Wrap(first_backtrace_point)
	return event
}

// Creates new domain chain on the backtrace
func (existing_event *Event) Contain(added_backtrace_point *trace.BacktracePoint) *Event {
	existing_event.Backtrace = append(existing_event.Backtrace, &BacktraceDomainChain{})
	dc := existing_event.Backtrace.Last(-1)
	*dc = append(*dc, added_backtrace_point)
	return existing_event
}

// Appends the event with additional backtrace point - just like the Reraise(), but without
// side effects
func (existing_event *Event) Wrap(added_backtrace_point *trace.BacktracePoint) *Event {
	dc := existing_event.Backtrace.Last(-1)
	*dc = append(*dc, added_backtrace_point)
	return existing_event
}
