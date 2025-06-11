package accumapi

import (
	"github.com/rusriver/configdb/v3"
	"github.com/rusriver/nutz/logger/loggerapi"
)

type ILogsAccumulator interface {
	loggerapi.ILogger
	SwitchOnTheOfflineAccumulation() loggerapi.ILogger
	FinishAccumulationAndFetch() (llines []*configdb.Root)
	RestartTheAccumulation()
}

/*
	Why separate? Because fucking idiotic Orwelian cycle import ban of Go language.
	Because another package, this depends on, really needs to use this logger too.

	Check the usage of it, the complexity:

		ol1 := Logger.NewSubLoggerInitChain("06fca0faf7").ILogger()
		var OfflineLogger accumapi.ILogsAccumulator = ol1.(*logger.Logger)
		OfflineLogger.SwitchOnTheOfflineAccumulation()

	But also all such unexpected changes break the established API contracts...
	The cycle imports ban was a REALLY BAD DECISION.
*/
