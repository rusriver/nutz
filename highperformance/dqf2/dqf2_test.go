package dqf2_test

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/rusriver/nutz/highperformance/dqf2"
	"github.com/stretchr/testify/require"
)

type Service01 struct {
	State int
	ChUp  chan *dqf2.Message[*Service01]
}

func (s *Service01) Run() {
	go func() {
		for m := range s.ChUp {
			if m.Func != nil {
				m.Func(m, s)
			}
		}
	}()
}

func Test_DQF2_01(t *testing.T) {
	s := &Service01{}
	s.ChUp = make(chan *dqf2.Message[*Service01], 100)
	s.Run()

	wg := &sync.WaitGroup{}
	for range 1000 {
		wg.Add(1)
		go func() {
			var x int
			mup := &dqf2.Message[*Service01]{Func: func(m *dqf2.Message[*Service01], s *Service01) {
				s.State++
				m.Tag = dqf2.Tag_OK
				x = s.State
			}}
			dqf2.SyncTransaction(s.ChUp, mup)
			fmt.Println("++", x, mup.Tag)
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
	require.Equal(t, 1000, x)
}

func Test_DQF2_02(t *testing.T) {
	s := &Service01{}
	s.ChUp = make(chan *dqf2.Message[*Service01], 20)
	s.Run()

	wg := &sync.WaitGroup{}
	for range 100 {
		wg.Add(1)
		go func() {
			var x int
			mup := &dqf2.Message[*Service01]{Func: func(m *dqf2.Message[*Service01], s *Service01) {
				s.State++
				m.Tag = dqf2.Tag_OK
				x = s.State
				time.Sleep(time.Millisecond * 50)
			}}
			dqf2.SyncTransaction(s.ChUp, mup)
			fmt.Println("++", x, mup.Tag)
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
	require.Equal(t, 100, x)
}

func Test_DQF2_03(t *testing.T) {
	s := &Service01{}
	s.ChUp = make(chan *dqf2.Message[*Service01], 20)
	s.Run()

	chdown := dqf2.NewChDown[*Service01]()

	wg := &sync.WaitGroup{}
	for range 100 {
		wg.Add(1)
		go func() {
			var x int
			mup := &dqf2.Message[*Service01]{Func: func(m *dqf2.Message[*Service01], s *Service01) {
				s.State++
				m.Tag = dqf2.Tag_OK
				x = s.State
				time.Sleep(time.Millisecond * 50)
			}}
			chdown.Transaction(s.ChUp, mup)
			<-chdown
			fmt.Println("++", x, mup.Tag)
			wg.Done()
		}()
	}
	wg.Wait()

	var x int
	m := &dqf2.Message[*Service01]{Func: func(m *dqf2.Message[*Service01], s *Service01) {
		m.Tag = dqf2.Tag_OK
		x = s.State
	}}
	chdown.Transaction(s.ChUp, m)
	<-chdown
	fmt.Println("==", x, m.Tag)
	require.Equal(t, 100, x)
}
