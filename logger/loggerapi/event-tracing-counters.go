package loggerapi

type EventTracingCounters []int

func (etc EventTracingCounters) Inc() {
	c := etc[len(etc)-1]
	c++
	if c > 0xFFFF {
		c = 0
	}
	etc[len(etc)-1] = c
}

func (etc EventTracingCounters) BranchSub() EventTracingCounters {
	etc2 := make(EventTracingCounters, len(etc)+1)
	copy(etc2, etc)
	return etc2
}
