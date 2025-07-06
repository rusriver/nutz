package dqf2

type ChDown[T any] chan *Message[T]

func NewChDown[T any]() (chDown ChDown[T]) {
	chDown = make(ChDown[T], 1)
	return chDown
}

// This can be used in select{}. See also the SyncTransaction() function.
// Note: please don't over-use this as a low-level primitive, it's best for high levels
// of abstraction code. For low-level, consider using the DQF pattern directly.
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
