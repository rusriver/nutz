package example1

import (
	"testing"

	"github.com/lithammer/shortuuid/v4"
	"github.com/rusriver/nutz/logger"
)

func Test_e3_01(t *testing.T) {

	ilogger := logger.New(nil).NSLIC(shortuuid.New()[:8]).ILogger()

	// get a zerolog logger
	zl := ilogger.(*logger.Logger).ZeroLogger

	// now log something
	zl.Info().Msg("Here's the zerologger output")
}
