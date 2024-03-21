package dqf

import (
	"sync"
	"time"
)

func SyncTransaction(ch chan func(), f func()) {
	wg := sync.WaitGroup{}
	wg.Add(1)
L:
	for {
		select {
		case ch <- func() {
			defer wg.Done()
			f()
		}:
			break L
		default:
			time.Sleep(time.Millisecond * 50)
		}
	}
	wg.Wait()
}
