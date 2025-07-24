package controlflow

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/rusriver/configdb/v3/deepcopy"
	"github.com/rusriver/nutz/loggeterr-api/trace"
)

type Event struct {
	RetryAttemptNb int `json:"-"`
	Forwardtrace   *trace.Forwardtrace
	Backtrace      FullBacktrace
}

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

var _ error = (*Event)(nil)

func (e *Event) Error() string {
	return e.String()
}

var _ fmt.Stringer = (*Event)(nil)

func (e *Event) String() string {
	return e.Serk()
}

func (e *Event) Serk() string {
	/*
		Example string:
			{: FORWARDTRACE {: ETC [1 2 3 4 5], TAGS [qq ww ee]}, BACKTRACE { [123123 "detail"] [123123 "detail"] [123123 "detail"], [123123 "detail"] [123123 "detail"]}}
	*/

	quote_f := func(s1 string) string {
		if strings.Contains(s1, " ") || strings.Contains(s1, "\t") || strings.Contains(s1, "\n") {
			return `"` + s1 + `"`
		} else {
			return s1
		}
	}

	br := strings.Builder{}
	br.WriteString("{: ")
	{
		if len(e.Forwardtrace.ETC) > 0 || len(e.Forwardtrace.Tags) > 0 {
			br.WriteString("FORWARDTRACE {: ")
			{
				if len(e.Forwardtrace.ETC) > 0 {
					br.WriteString("ETC [")
					{
						l := len(e.Forwardtrace.ETC) - 1
						for i, v := range e.Forwardtrace.ETC {
							br.WriteString(strconv.Itoa(v))
							if i < l {
								br.WriteString(" ")
							}
						}
					}
					br.WriteString("], ")
				}
				//
				if len(e.Forwardtrace.Tags) > 0 {
					br.WriteString("TAGS [")
					{
						l := len(e.Forwardtrace.Tags) - 1
						for i, v := range e.Forwardtrace.Tags {
							br.WriteString(quote_f(v))
							if i < l {
								br.WriteString(" ")
							}
						}
					}
					br.WriteString("]")
				}
			}
			br.WriteString("}, ")
		}
		//
		br.WriteString("BACKTRACE { ")
		{
			l1 := len(e.Backtrace) - 1
			for i1, dc := range e.Backtrace {

				l2 := len(*dc) - 1
				for i2, bp := range *dc {
					br.WriteString("[")
					br.WriteString(quote_f(bp.Id))
					br.WriteString(" ")
					br.WriteString(quote_f(fmt.Sprintf("%v", bp.Detail)))
					br.WriteString("]")
					if i2 < l2 {
						br.WriteString(" ")
					}
				}

				if i1 < l1 {
					br.WriteString(", ")
				}
			}
		}
		br.WriteString("}")
	}
	br.WriteString("}")

	return br.String()
}

func (e *Event) JSON() []byte {
	bb, _ := json.Marshal(e)
	return bb
}

func Try(f func() (err any)) (e *Event) {
	defer func() {
		if x := recover(); x != nil {
			switch xv := x.(type) {
			case *trace.BacktracePoint:
				e = NewEvent(nil, xv)
			case *Event:
				e = xv
			case error:
				e = NewEvent(nil, &trace.BacktracePoint{"24ba97c201fe1", xv.Error()})
			default:
				e = NewEvent(nil, trace.BacktracePointf("24ba97c201fe2", "%v", xv))
			}
		}
	}()
	x := f()
	if x != nil {
		switch xv := x.(type) {
		case *trace.BacktracePoint:
			e = NewEvent(nil, xv)
		case *Event:
			e = xv
		case error:
			e = NewEvent(nil, &trace.BacktracePoint{"24ba97c201fe3", xv.Error()})
		default:
			e = NewEvent(nil, trace.BacktracePointf("24ba97c201fe4", "%v", xv))
		}
	}
	return
}

func (e *Event) Catch(f func(e *Event)) *Event {
	if e != nil {
		f(e)
	}
	return e
}

func (e *Event) OK(f func()) *Event {
	if e == nil {
		f()
	}
	return e
}

// Allows to try and retry an operation several times.
func Retry(f_do func() (err any), f_failure func(e *Event) (end bool)) (e *Event) {
	attemptNb := 0
	for {
		attemptNb++
		e = Try(func() (err any) {
			err = f_do()
			return
		})
		if e == nil {
			break
		} else {
			e.RetryAttemptNb = attemptNb
			end := f_failure(e)
			if end {
				break
			}
		}
	}
	return
}

func (existing_event *Event) Reraise(added_backtrace_point *trace.BacktracePoint) {
	if existing_event != nil {
		if added_backtrace_point != nil {
			existing_event.Wrap(added_backtrace_point)
		}
		panic(existing_event)
	}
}

/*
Example usage:

	controlflow.RaiseCause(tracer.GetForwardtrace(), trace.BacktracePointf("123123", "the cause"))
	controlflow.RaiseCause(nil, trace.BacktracePointf("123123", "the cause"))

Sometimes, some function want to raise the cause, but it's too tedious to pass the tracer in it;
in this case, the function may raise it with nil forwardtrace (and it's okay), and then the user
of it will catch the event, and set actual forwardtrace directly. Example, inside such a function:

	controlflow.RaiseCause(nil, trace.BacktracePointf("123123", err.Error()))

And then at first interception point:

	event.Forwardtrace = tracer.GetForwardtrace()
*/
func RaiseCause(forward_trace *trace.Forwardtrace, first_backtrace_point *trace.BacktracePoint) {
	event := NewEvent(forward_trace, first_backtrace_point)
	panic(event)
}

// Example usage: controlflow.NewEvent(tracer.GetForwardtrace(), trace.BacktracePointf("123123", "the cause"))
// controlflow.NewEvent(nil, trace.BacktracePointf("123123", "the cause"))
func NewEvent(forward_trace *trace.Forwardtrace, first_backtrace_point *trace.BacktracePoint) *Event {
	if forward_trace == nil {
		forward_trace = &trace.Forwardtrace{}
	}
	event := &Event{
		Forwardtrace: deepcopy.Copy(forward_trace).(*trace.Forwardtrace),
		Backtrace: FullBacktrace{
			&BacktraceDomainChain{},
		},
	}
	event.Wrap(first_backtrace_point)
	return event
}

// Creates new domain chain on the backtrace
func (existing_event *Event) ContainAndWrap(added_backtrace_point *trace.BacktracePoint) *Event {
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
