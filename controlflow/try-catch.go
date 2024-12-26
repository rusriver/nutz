package controlflow

type Exception struct {
	X interface{}
}

func Try(f func()) (e *Exception) {
	defer func() {
		if x := recover(); x != nil {
			e = &Exception{
				X: x,
			}
		}
	}()
	f()
	return
}

func (e *Exception) Catch(f func(e *Exception)) {
	if e != nil {
		f(e)
	}
}
