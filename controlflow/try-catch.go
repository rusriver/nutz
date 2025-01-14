package controlflow

import "fmt"

type Exception struct {
	X   interface{}
	Err error
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

func (e *Exception) Catch(f func(e *Exception)) {
	if e != nil {
		f(e)
	}
}
