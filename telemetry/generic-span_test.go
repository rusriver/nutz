package telemetry_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/rusriver/nutz/telemetry"
)

type MySpanAttrs struct {
	HttpCode int
	Path     string
}

func Test_span_001(t *testing.T) {
	spanCreator := telemetry.NewSpanCreator(func(c *telemetry.SpanCreator[*MySpanAttrs]) {
		c.StarterFunc = func(s *telemetry.Span[*MySpanAttrs]) {
			fmt.Println("++ starter")
			s.T0 = time.Now()
		}
		c.FinisherFunc = func(s *telemetry.Span[*MySpanAttrs]) {
			fmt.Printf("++finisher %v %+v", s.T1.Sub(s.T0), s.UserAttrs)
		}
		c.ErrorFunc = func(s *telemetry.Span[*MySpanAttrs], err error) {
			fmt.Printf("++error %v %+v", err, s.UserAttrs)
		}
	})

	span := spanCreator.NewSpan(func(s *telemetry.Span[*MySpanAttrs]) {
		s.UserAttrs = &MySpanAttrs{}
	})

	span.UserAttrs.HttpCode = 404
	span.UserAttrs.Path = "/hello"

	// these are equivalent ways
	span.T0 = time.Now()
	span.Start()

	time.Sleep(time.Millisecond * 350)

	span.T1 = time.Now()

	time.Sleep(time.Second * 2)

	span.Finish()

	/*
		++ starter
		++finisher 363.7091ms &{HttpCode:404 Path:/hello}
	*/
}
