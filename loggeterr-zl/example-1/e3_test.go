package example1

import (
	"testing"

	"github.com/lithammer/shortuuid/v4"
	loggeterrzl "github.com/rusriver/nutz/loggeterr-zl"
)

func Test_e3_01(t *testing.T) {

	itracer := loggeterrzl.New(nil).NSTIC(shortuuid.New()[:8]).ITracer()

	// get a zerolog logger
	zl := itracer.(*loggeterrzl.Tracer).ZeroLoggerPtr

	// now log something
	zl.Info().Msg("Here's the zerologger output")
}
