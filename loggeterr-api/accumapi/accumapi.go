package accumapi

import (
	"github.com/rusriver/configdb/v3"
	loggeterrapi "github.com/rusriver/nutz/loggeterr-api"
)

type ILogsAccumulator interface {
	loggeterrapi.ITracer
	SwitchOnTheOfflineAccumulation() loggeterrapi.ITracer
	FinishAccumulationAndFetch() (llines []*configdb.Root)
	RestartTheAccumulation()
}

/*
	Why separate? Because fucking idiotic Orwelian cycle import ban of Go language.
	Because another package, this depends on, really needs to use this logger too.

	Check the usage of it, the complexity:

		ot1 := tracer.NewSubTracerInitChain("06fca0faf7").ITracer()
		var OfflineTracer accumapi.ILogsAccumulator = ot1.(*loggeterr.Tracer)
		OfflineTracer.SwitchOnTheOfflineAccumulation()

	But also all such unexpected changes break the established API contracts...
	The cycle imports ban was a REALLY BAD DECISION.
*/
