package controlflow

type Exception struct {
	X interface{}
}

func Try(f func() (err error)) (e *Exception) {
	defer func() {
		if x := recover(); x != nil {
			e = &Exception{
				X: x,
			}
		}
	}()
	err := f()
	if err != nil {
		e = &Exception{
			X: err,
		}
	}
	return
}

func (e *Exception) Catch(f func(e *Exception)) {
	if e != nil {
		f(e)
	}
}
