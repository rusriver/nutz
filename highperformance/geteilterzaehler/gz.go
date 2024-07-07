package geteilterzaehler

import (
	"context"
	"sync"
	"time"
)

type Int struct {
	Breadth   uint16
	schuetzen uint16 // the race shuttle
	scherben  []*Scherbe
}

type Scherbe struct {
	V  float64
	mu sync.Mutex
}

func NewInt(ff ...func(gz *Int)) (gz *Int) {
	gz = &Int{
		Breadth: 16,
	}
	for _, f := range ff {
		f(gz)
	}
	gz.scherben = make([]*Scherbe, gz.Breadth)
	for i := 0; i < int(gz.Breadth); i++ {
		gz.scherben[i] = &Scherbe{V: 0}
	}
	return gz
}

func (gz *Int) ApplyValue(extraEntropy uint16, f func(s *Scherbe)) {
	// intentionally subject to the race
	si := gz.schuetzen
	si += 1 + extraEntropy%gz.Breadth
	if si >= gz.Breadth {
		si = 0
	}
	gz.schuetzen = si

	s := gz.scherben[si]
	s.mu.Lock()
	f(s)
	s.mu.Unlock()
}

func (gz *Int) Reduce(f func(s *Scherbe)) {
	for _, s := range gz.scherben {
		s.mu.Lock()
		f(s)
		s.mu.Unlock()
	}
	return
}

// A backgroud goroutine for constant-rate reduction process
func (gz *Int) BackgroundReducer(ctx context.Context, tickPeriod time.Duration, batchPerTick int, f func(s *Scherbe, nextRound bool)) {
	tick := time.NewTicker(tickPeriod)
	i := 0
	for {
		select {
		case <-ctx.Done():
			return
		case <-tick.C:
			for bi := 0; bi < batchPerTick; bi++ {
				if i == 0 {
					f(gz.scherben[i], true)
				} else {
					f(gz.scherben[i], false)
				}
				i++
				if i == int(gz.Breadth) {
					i = 0
				}
			}
		}
	}
}
