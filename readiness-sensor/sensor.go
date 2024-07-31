package readinesssensor

import (
	"context"
	"time"

	"github.com/rusriver/nutz/highperformance/dqf"
	"github.com/rusriver/ttlcache/v3"
)

type Sensor struct {
	Registry *ttlcache.Cache[string, bool]
	TTL      time.Duration
	Context  context.Context

	DQF     chan func(s *Sensor)
	DQF_Len int
}

func New(of ...func(c *Sensor)) (s *Sensor, err error) {
	s = &Sensor{
		TTL:     time.Duration(0),
		Context: context.Background(),
		DQF_Len: 124,
	}
	for _, of := range of {
		of(s)
	}
	s.DQF = make(chan func(s *Sensor), s.DQF_Len)

	if s.TTL > 0 {
		s.Registry = ttlcache.New(ttlcache.WithTTL[string, bool](s.TTL))
		go s.Registry.Start()
	} else {
		s.Registry = ttlcache.New[string, bool]()
	}

	// main worker
	go func() {
		for {
			select {
			case <-s.Context.Done():
				s.Registry.Stop()
				return
			case dqf := <-s.DQF:
				dqf(s)
			}
		}
	}()

	return
}

func (s *Sensor) Report(w string) {
	s.DQF <- func(s *Sensor) {
		s.Registry.Set(w, true)
	}
}

func (s *Sensor) GetAllList() (all []string) {
	dqf.SyncTransaction(s.DQF, func(s *Sensor) {
		all = s.Registry.Keys()
	})
	return
}

func (s *Sensor) MatchTheProfile(prof []string) (ok bool, match []string, missing []string) {
	all := s.GetAllList()

	allMap := make(map[string]bool, len(all))
	for _, k := range all {
		allMap[k] = true
	}

	for _, k := range prof {
		if allMap[k] {
			match = append(match, k)
		} else {
			missing = append(missing, k)
		}
	}

	if len(prof) == len(match) {
		ok = true
	} else {
		ok = false
	}
	return
}
