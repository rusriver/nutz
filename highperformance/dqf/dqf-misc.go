package dqf

import "sync"

func SyncTransaction(ch chan func(), f func()) {
	wg := sync.WaitGroup{}
	wg.Add(1)
	ch <- func() {
		defer wg.Done()
		f()
	}
	wg.Wait()
}
