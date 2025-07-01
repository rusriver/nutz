package dqf2_test

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/rusriver/nutz/highperformance/dqf2"
	loggeterrzl "github.com/rusriver/nutz/loggeterr-zl"
	"github.com/stretchr/testify/require"
)

func Test_Chopper_01(t *testing.T) {

	// possible production profile
	au_modes := []dqf2.Chopper_Mode{
		{
			Period:              3000 * time.Millisecond,
			ThresholdPerc:       90,
			ChopNPerc:           20,
			NChopsToUpgrade:     3,
			NNoChopsToDowngrade: 8,
		},
		{
			Period:              800 * time.Millisecond,
			ThresholdPerc:       80,
			ChopNPerc:           10,
			NChopsToUpgrade:     3,
			NNoChopsToDowngrade: 16,
		},
		{
			Period:              200 * time.Millisecond,
			ThresholdPerc:       80,
			ChopNPerc:           10,
			NChopsToUpgrade:     3,
			NNoChopsToDowngrade: 32,
		},
		{
			Period:              50 * time.Millisecond,
			ThresholdPerc:       80,
			ChopNPerc:           10,
			NChopsToUpgrade:     3,
			NNoChopsToDowngrade: 64,
		},
	}

	// for testing specifically
	au_modes = []dqf2.Chopper_Mode{
		{
			Period:              800 * time.Millisecond,
			ThresholdPerc:       80,
			ChopNPerc:           10,
			NChopsToUpgrade:     3,
			NNoChopsToDowngrade: 16,
		},
		{
			Period:              200 * time.Millisecond,
			ThresholdPerc:       80,
			ChopNPerc:           10,
			NChopsToUpgrade:     3,
			NNoChopsToDowngrade: 32,
		},
		{
			Period:              50 * time.Millisecond,
			ThresholdPerc:       80,
			ChopNPerc:           10,
			NChopsToUpgrade:     3,
			NNoChopsToDowngrade: 64,
		},
	}

	ctx := context.Background()
	tracer := loggeterrzl.New(nil)

	s := &Service01{}
	s.ChUp = make(chan *dqf2.Message[*Service01], 20)
	s.Run()

	auTracer := tracer.NSTIC("au01").ITracer()
	dqf2.Chopper_Start(ctx, s.ChUp, func(au *dqf2.Chopper_Options) {
		au.Modes = au_modes
		au.IfChoppedFunc = func(n_chopped int) {
			auTracer.InfoEvent("55").Int("n-chopped", n_chopped).Send()
		}
		au.UpgradeFunc = func(from_mode, to_mode int) {
			auTracer.InfoEvent("88+").SendMsgf("upgrade %v -> %v", from_mode, to_mode)
		}
		au.DowngradeFunc = func(from_mode, to_mode int) {
			auTracer.InfoEvent("88-").SendMsgf("downgrade %v -> %v", from_mode, to_mode)
		}
	})

	wg := &sync.WaitGroup{}
	for range 500 {
		wg.Add(1)
		go func() {
			var x int
			m := &dqf2.Message[*Service01]{Func: func(m *dqf2.Message[*Service01], s *Service01) {
				s.State++
				m.Tag = dqf2.Tag_OK
				x = s.State
				time.Sleep(time.Millisecond * 100)
			}}
			dqf2.SyncTransaction(s.ChUp, m)
			fmt.Println("++", x, m.Tag)
			wg.Done()
		}()
	}
	wg.Wait()

	var x int
	m := &dqf2.Message[*Service01]{Func: func(m *dqf2.Message[*Service01], s *Service01) {
		m.Tag = dqf2.Tag_OK
		x = s.State
	}}
	dqf2.SyncTransaction(s.ChUp, m)
	fmt.Println("==", x, m.Tag)
	require.InDelta(t, 92, x, 5)

	wg = &sync.WaitGroup{}
	for range 100 {
		wg.Add(1)
		go func() {
			var x int
			m := &dqf2.Message[*Service01]{Func: func(m *dqf2.Message[*Service01], s *Service01) {
				s.State++
				m.Tag = dqf2.Tag_OK
				x = s.State
				time.Sleep(time.Millisecond * 100)
			}}
			dqf2.SyncTransaction(s.ChUp, m)
			fmt.Println("++", x, m.Tag)
			wg.Done()
		}()
		time.Sleep(time.Millisecond * 100)
	}
	wg.Wait()

	m = &dqf2.Message[*Service01]{Func: func(m *dqf2.Message[*Service01], s *Service01) {
		m.Tag = dqf2.Tag_OK
		x = s.State
	}}
	dqf2.SyncTransaction(s.ChUp, m)
	fmt.Println("==", x, m.Tag)
	require.InDelta(t, 192, x, 5)
}
