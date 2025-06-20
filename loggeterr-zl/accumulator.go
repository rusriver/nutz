package loggeterrzl

import (
	"sync"

	"github.com/rusriver/configdb/v3"
	loggeterrapi "github.com/rusriver/nutz/loggeterr-api"
)

type Accumulator struct {
	logsAccumulated []*configdb.Root
	chUp            chan []byte
	wg              sync.WaitGroup
}

func NewAccumulator(tracer loggeterrapi.ITracer) *Accumulator {
	a := &Accumulator{
		chUp: make(chan []byte, 100),
	}

	a.wg.Add(1)

	go func() {
		for lline := range a.chUp {
			db := (&configdb.InitContext{}).WithTracer(tracer).FromBytes(lline).Load().DontPanic()
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
