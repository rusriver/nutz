package dqf2

type ChDown[T any] chan *Message[T]

func NewChDown[T any]() (chDown ChDown[T]) {
	chDown = make(ChDown[T], 1)
	return chDown
}

// This can be used in select{}. See also the SyncTransaction() function.
func (chDown ChDown[T]) Transaction(chUp chan *Message[T], m *Message[T]) {
	select {
	case <-chDown:
	default:
	}
	m.chDown = chDown
	f := m.Func
	m.Func = func(m *Message[T], s T) {
		f(m, s)
		m.chDown <- m
	}
	chUp <- m
}
