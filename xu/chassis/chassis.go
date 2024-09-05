package chassis

import (
	"math/rand/v2"
	"os"

	"github.com/rusriver/config/v2"
	"github.com/rusriver/nutz/logger"
	readinesssensor "github.com/rusriver/nutz/readiness-sensor"
	promethabs "github.com/rusriver/nutz/telemetry/prometh-abs"
	valvidaabs "github.com/rusriver/nutz/valvida/valvida-abs"
)

type Chassis struct {
	Logger  logger.ILogger
	Metrics promethabs.IMetricsService

	// If not nil, the chassis will try to go by the structural config path,
	// else by the free-form config path. See the docs for details.
	StructConfig any

	// This only gets set if you go the free-form config path and set StructConfig to nil
	ConfigSource *config.Source

	Inventory       *readinesssensor.Sensor
	DiagnosticFlags map[string]bool

	Logging struct {
		RL valvidaabs.IRatelimService
		/*
			TODO: Please see an arec on how to implement this.
		*/

		CardinalityLimit struct {
			Table              map[string]struct{}
			Lim                int
			СhCardinalityLimUp chan string
		}
	}

	USD any

	ChShutdown   chan os.Signal
	PseudoRandom *rand.Rand
}

func New(fo ...func(c *Chassis)) *Chassis {
	c := &Chassis{}

	return c
}
