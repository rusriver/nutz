package controlflow

import (
	"fmt"
	"strings"
)

type Exception struct {
	X             interface{}
	TheCauseError error
	AttemptNb     int
	XChain        []any
}

var _ error = (*Exception)(nil)

func (e *Exception) Error() string {
	return e.String()
}

var _ fmt.Stringer = (*Exception)(nil)

func (e *Exception) String() string {
	br := strings.Builder{}
	for i := len(e.XChain) - 1; i >= 0; i-- {
		x := e.XChain[i]
		switch xv := x.(type) {
		case error:
			br.WriteString(xv.Error())
			br.WriteString("; ")
		case string:
			br.WriteString(xv)
			br.WriteString("; ")
		default:
			br.WriteString(fmt.Sprintf("%v; ", xv))
		}
	}
	br.WriteString(
		e.TheCauseError.Error(),
	)
	return br.String()
}

func Try(f func() (err error)) (e *Exception) {
	defer func() {
		if x := recover(); x != nil {
			switch xv := x.(type) {
			case *Exception:
				e = xv
			case error:
				e = &Exception{
					X:             x,
					TheCauseError: xv,
				}
			default:
				e = &Exception{
					X:             x,
					TheCauseError: fmt.Errorf("%v", xv),
				}
			}
		}
	}()
	err := f()
	if err != nil {
		e = &Exception{
			X:             err,
			TheCauseError: err,
		}
	}
	return
}

func (e *Exception) Catch(f func(e *Exception)) *Exception {
	if e != nil {
		f(e)
	}
	return e
}

func (e *Exception) OK(f func(e *Exception)) *Exception {
	if e == nil {
		f(e)
	}
	return e
}

// Allows to try and retry an operation several times.
func Retry(f func() (err error), f2 func(e *Exception) (end bool)) (e *Exception) {
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
			e.AttemptNb = attemptNb
			end := f2(e)
			if end {
				break
			}
		}
	}
	return
}

// If you added context's string representation is longer than 15 bytes, then
// maybe you are using it wrongly. In 99% of cases, the added context should be a string.
func Raise2(added_context any, reason any) {
	switch xv := reason.(type) {
	case *Exception:
		xv.XChain = append(xv.XChain, added_context)
		panic(xv)
	case error:
		x2 := &Exception{
			X:             xv,
			TheCauseError: xv,
		}
		x2.XChain = append(x2.XChain, added_context)
		panic(x2)
	default:
		x2 := &Exception{
			X:             xv,
			TheCauseError: fmt.Errorf("%v", xv),
		}
		x2.XChain = append(x2.XChain, added_context)
		panic(x2)
	}
}

// You better not use Raise1() to re-raise exceptions - use Raise2() instead.
// This one is best intended for cases where you create a new exception.
func Raise1(reason any) {
	switch xv := reason.(type) {
	case *Exception:
		panic(xv)
	case error:
		x2 := &Exception{
			X:             xv,
			TheCauseError: xv,
		}
		panic(x2)
	default:
		x2 := &Exception{
			X:             xv,
			TheCauseError: fmt.Errorf("%v", xv),
		}
		panic(x2)
	}
}
