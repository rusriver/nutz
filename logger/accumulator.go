package logger

import (
	"sync"

	"github.com/rusriver/configdb/v3"
)

type Accumulator struct {
	logsAccumulated []*configdb.Root
	chUp            chan []byte
	wg              sync.WaitGroup
}

func NewAccumulator() *Accumulator {
	a := &Accumulator{
		chUp: make(chan []byte, 100),
	}

	a.wg.Add(1)

	go func() {
		for lline := range a.chUp {
			db := (&configdb.InitContext{}).FromBytes(lline).Load()
			a.logsAccumulated = append(a.logsAccumulated, db)
		}
		a.wg.Done()
	}()

	return a
}

func (a *Accumulator) FinishAccumulationAndFetch() []*configdb.Root {
	close(a.chUp)
	a.wg.Wait()
	return a.logsAccumulated
}
