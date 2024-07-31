package timmer

import "time"

// This is a wrapper around Go std Timer, with a saner API.
type Timmer struct {
	Timer *time.Timer
	C     <-chan time.Time
}

func New() *Timmer {
	t := &Timmer{}
	t.Timer = time.NewTimer(1)
	t.Stop()
	t.C = t.Timer.C
	return t
}

func (t *Timmer) Stop() {
	t.Timer.Stop()
	select {
	case <-t.Timer.C:
	default:
	}
}

func (t *Timmer) Restart(d time.Duration) {
	t.Stop()
	t.Timer.Reset(d)
}
