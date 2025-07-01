package dqf2

import (
	"context"
	"sync"
	"time"

	"github.com/rusriver/nutz/controlflow"
	"github.com/rusriver/nutz/loggeterr-api/trace"
)

type Message[T any] struct {
	Tag  int // this is used, among other things, by the chopper to mark failed deliveries
	wg   sync.WaitGroup
	Func func(m *Message[T], s T)
}

const (
	Tag_Fresh   = 0   // the default for new messages, or must be reset by clients
	Tag_OK      = 1   // set by server only
	Tag_Dropped = 177 // set by chopper only
)

func SyncTransaction[T any](chUp chan *Message[T], m *Message[T]) {
	m.wg.Add(1)
	f := m.Func
	m.Func = func(m *Message[T], s T) {
		defer m.wg.Done()
		f(m, s)
	}
	chUp <- m
	m.wg.Wait()
}

type Chopper_Options struct {
	Modes         []Chopper_Mode // 0 is slowest, -1 is fastest
	IfChoppedFunc func(n_chopped int)
	UpgradeFunc   func(from_mode, to_mode int)
	DowngradeFunc func(from_mode, to_mode int)
}

type Chopper_Mode struct {
	Period              time.Duration
	ThresholdPerc       float32 // threshold perc of full capacity, action happens if above it
	ChopNPerc           float32 // how much to chop per one action, perc of full capacity
	NChopsToUpgrade     int     // e.g. after 5 successive chops, upgrade to faster mode
	NNoChopsToDowngrade int     // e.g. after 5 successive no-chops, downgrade to slower mode
}

func Chopper_Start[T any](
	ctx context.Context,
	chUp chan *Message[T],
	of ...func(*Chopper_Options),
) {
	opts := &Chopper_Options{
		// the defaults
	}
	for _, f := range of {
		f(opts)
	}

	if len(opts.Modes) == 0 {
		controlflow.RaiseCause(nil, &trace.BacktracePoint{"0197c5e9eb3249560341", "zero number of modes"})
	}

	go func() {
		tickers := []*time.Ticker{}
		for _, mode := range opts.Modes {
			tickers = append(tickers, time.NewTicker(mode.Period))
		}

		current_mode_i := 0 // slowest by default
		a := 0              // + chops, - no-chops

		for {
			current_mode := opts.Modes[current_mode_i]
			select {
			case <-ctx.Done():
				return
			case <-tickers[current_mode_i].C:
				L := float32(len(chUp))
				C := float32(cap(chUp))
				over := L/C*100 - current_mode.ThresholdPerc
				n := int(over / 100 * C)
				if n > 0 {
					a++
					for range n {
						m := <-chUp
						m.Tag = Tag_Dropped
						m.wg.Done()
					}
					if opts.IfChoppedFunc != nil {
						opts.IfChoppedFunc(n)
					}
				} else {
					a--
				}
				switch {
				case a >= current_mode.NChopsToUpgrade:
					new_mode_i := current_mode_i + 1
					if new_mode_i >= len(opts.Modes) {
						new_mode_i--
					}
					if opts.UpgradeFunc != nil {
						opts.UpgradeFunc(current_mode_i, new_mode_i)
					}
					current_mode_i = new_mode_i
					a = 0
				case a <= 0-current_mode.NNoChopsToDowngrade:
					new_mode_i := current_mode_i - 1
					if new_mode_i < 0 {
						new_mode_i = 0
					}
					if opts.DowngradeFunc != nil {
						opts.DowngradeFunc(current_mode_i, new_mode_i)
					}
					current_mode_i = new_mode_i
					a = 0
				}
			}
		}
	}()
}
