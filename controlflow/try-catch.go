package controlflow

import "fmt"

type Exception struct {
	AttemptNb int
	X         interface{}
	Err       error
}

func Try(f func() (err error)) (e *Exception) {
	defer func() {
		if x := recover(); x != nil {
			e = &Exception{
				X:   x,
				Err: fmt.Errorf("%v", x),
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
					e = &Exception{
						X:   x,
						Err: fmt.Errorf("%v", x),
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
