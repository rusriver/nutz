package dqf

import (
	"sync"
	"time"
)

func SyncTransaction[T any](ch chan func(T), f func(T), retry ...time.Duration) {
	var rd time.Duration
	if len(retry) > 0 {
		rd = retry[0]
	} else {
		rd = time.Millisecond * 50
	}
	wg := sync.WaitGroup{}
	wg.Add(1)
L:
	for {
		select {
		case ch <- func(v T) {
			defer wg.Done()
			f(v)
		}:
			break L
		default:
			time.Sleep(rd)
		}
	}
	wg.Wait()
}
