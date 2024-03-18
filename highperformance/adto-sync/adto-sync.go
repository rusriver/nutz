package adtosync

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"time"
)

type ADTOSync struct {
	mu          sync.Mutex
	ctx         context.Context
	cancelCtx   context.CancelFunc
	err         error
	once        FastAsyncOnce
	Invalidated bool
}

func (ds *ADTOSync) SealWithTimeout(ctxUp context.Context, to time.Duration) {
	ds.mu.Lock()
	ds.ctx, ds.cancelCtx = context.WithCancel(ctxUp)
	go func() {
		select {
		case <-time.After(to):
			// timeout path
			ds.once.Do(func() {
				ds.Invalidated = true
				ds.err = errors.New("f633da35 Timeout")
				ds.mu.Unlock()
			})
			return
		case <-ds.ctx.Done():
			return
		}
	}()
}

func (ds *ADTOSync) Wait() (err error) {
	// wait for anyone unlocking it
	ds.mu.Lock()
	ds.mu.Unlock()
	return ds.err
}

func (ds *ADTOSync) Unseal() {
	// normal path
	ds.cancelCtx()
	ds.once.Do(func() {
		ds.mu.Unlock()
	})
}

/*
	func (o *Once) Do(f func()) {
		// Note: Here is an incorrect implementation of Do:
		//
		//	if atomic.CompareAndSwapUint32(&o.done, 0, 1) {
		//		f()
		//	}
		//
		// Do guarantees that when it returns, f has finished.
		// This implementation would not implement that guarantee:
		// given two simultaneous calls, the winner of the cas would
		// call f, and the second would return immediately, without
		// waiting for the first's call to f to complete.
		// This is why the slow path falls back to a mutex, and why
		// the atomic.StoreUint32 must be delayed until after f returns.

p.s. we don't need that guarantee here, therefore need a CAS specifically
*/
type FastAsyncOnce struct {
	done uint32
}

func (o *FastAsyncOnce) Do(f func()) {
	if atomic.CompareAndSwapUint32(&o.done, 0, 1) {
		f()
	}
}
