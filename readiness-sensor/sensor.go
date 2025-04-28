package readinesssensor

import (
	"context"
	"time"

	"github.com/rusriver/nutz/highperformance/dqf"
	"github.com/rusriver/nutz/logger"
	"github.com/rusriver/ttlcache/v3"
)

type Sensor struct {
	Registry *ttlcache.Cache[string, uint8]
	TTL      time.Duration
	Context  context.Context
	Logger   logger.ILogger
	F        func(s *Sensor)

	// Called every time the list changes - when anything is added, or removed by TTL
	ChangeFunc func(s *Sensor, what string)

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
		s.Registry = ttlcache.New(ttlcache.WithTTL[string, uint8](s.TTL))
		go s.Registry.Start()
	} else {
		s.Registry = ttlcache.New[string, uint8]()
	}

	if s.ChangeFunc != nil {
		s.Registry.OnInsertion(func(ctx context.Context, i *ttlcache.Item[string, uint8]) {
			s.ChangeFunc(s, i.Key())
		})
		s.Registry.OnEviction(func(ctx context.Context, er ttlcache.EvictionReason, i *ttlcache.Item[string, uint8]) {
			s.ChangeFunc(s, i.Key())
		})
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

func (s *Sensor) Report(what string, status uint8) {
	s.DQF <- func(s *Sensor) {
		s.Registry.Set(what, status)
		if s.F != nil {
			s.F(s)
		}
	}
}

func (s *Sensor) NewProbe(name string) *Probe {
	probe := &Probe{
		Sensor: s,
		Name:   name,
	}
	probe.Sensor.Report(name, SubsystemStatus_Fresh)
	return probe
}

// DO NOT use from within the reactor func - will deadlock
func (s *Sensor) GetAllList() (all []string) {
	dqf.SyncTransaction(s.DQF, func(s *Sensor) {
		all = s.Registry.Keys()
	})
	return
}

// DO NOT use from within the reactor func - will deadlock
func (s *Sensor) GetAllItems() (all map[string]uint8) {
	all = make(map[string]uint8)
	dqf.SyncTransaction(s.DQF, func(s *Sensor) {
		for k, item := range s.Registry.Items() {
			all[k] = item.Value()
		}
	})
	return
}

// DO NOT use from within the reactor func - will deadlock
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

type Status struct {
	Fresh   []string
	Ready   []string
	Ceded   []string
	AllDone bool
}

func (s *Sensor) Status(all_items map[string]uint8) (st *Status) {
	st = &Status{
		AllDone: true,
	}
	for ss, status := range all_items {
		switch status {
		case SubsystemStatus_Fresh:
			st.Fresh = append(st.Fresh, ss)
			st.AllDone = false
		case SubsystemStatus_Ready:
			st.Ready = append(st.Ready, ss)
		case SubsystemStatus_Ceded:
			st.Ceded = append(st.Ceded, ss)
		}
	}
	return
}

func (s *Sensor) StatusInternal() (st *Status) {
	all_items := map[string]uint8{}
	for k, item := range s.Registry.Items() {
		all_items[k] = item.Value()
	}
	return s.Status(all_items)
}
