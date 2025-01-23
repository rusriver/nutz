package controlflow

import (
	"fmt"
	"strings"
)

type Exception struct {
	X         interface{}
	Err       error
	AttemptNb int
	XChain    []any
}

func (e *Exception) String() string {
	br := strings.Builder{}
	for i := len(e.XChain) - 1; i >= 0; i-- {
		x := e.XChain[i]
		switch xv := x.(type) {
		case error:
			br.WriteString(xv.Error())
			br.WriteString(" ")
		case string:
			br.WriteString(xv)
			br.WriteString(" ")
		default:
			br.WriteString(fmt.Sprintf("%v ", xv))
		}
	}
	br.WriteString(
		e.Err.Error(),
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
					X:   x,
					Err: xv,
				}
			default:
				e = &Exception{
					X:   x,
					Err: fmt.Errorf("%v", xv),
				}
			}
		}
	}()
	err := f()
	if err != nil {
		e = &Exception{
			X:   err,
			Err: err,
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
		e = nil
		func() {
			defer func() {
				if x := recover(); x != nil {
					switch xv := x.(type) {
					case *Exception:
						e = xv
					case error:
						e = &Exception{
							X:   x,
							Err: xv,
						}
					default:
						e = &Exception{
							X:   x,
							Err: fmt.Errorf("%v", xv),
						}
					}
				}
			}()
			err := f()
			if err != nil {
				e = &Exception{
					X:   err,
					Err: err,
				}
			}
		}()
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

func Raise(x any) {
	panic(x)
}

func Raise2(added_context any, x any) {
	switch xv := x.(type) {
	case *Exception:
		xv.XChain = append(xv.XChain, added_context)
		panic(xv)
	case error:
		x2 := &Exception{
			X:   xv,
			Err: xv,
		}
		panic(x2)
	default:
		x2 := &Exception{
			X:   xv,
			Err: fmt.Errorf("%v", xv),
		}
		panic(x2)
	}
}
