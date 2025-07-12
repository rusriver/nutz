package example1_test

import (
	"errors"
	"fmt"
	"sync"
	"testing"

	"github.com/lithammer/shortuuid/v4"
	"github.com/rusriver/nutz/controlflow"
	loggeterrapi "github.com/rusriver/nutz/loggeterr-api"
	"github.com/rusriver/nutz/loggeterr-api/trace"
	loggeterrzl "github.com/rusriver/nutz/loggeterr-zl"
)

func Test_e1_01(t *testing.T) {

	tracer := loggeterrzl.New(nil).NSTIC(shortuuid.New()[:8]).ITracer()

	var Tracer loggeterrapi.ITracer

	Tracer = tracer

	Tracer.LogEvent("d3740cd9").Title("some title 1").Send()
	Tracer.LogEvent("05c1e865").Title("some title 1").Send()

	Tracer = tracer.NewSubTracerInitChain("SUBSYSTEM-a24487").
		Str("xu", "my-xu").
		Str("ver", "1.2.3").
		ITracer()

	Tracer.LogEvent("111b7de0").Title("some title 1").Send()
	Tracer.LogEvent("e3466c79").Title("some title 1").Send()

	Tracer.SpanBegin("SPAN-3a7ce735")
	Tracer.LogEvent("32").SendMsg("hello")
	Tracer.LogEvent("33").SendMsg("again")
	Tracer.SpanEnd()

	Tracer.LogEvent("66afccbd").Title("some title 1").Send()

	Tracer = Tracer.NewSubTracerInitChain("59fc7").Str("some", "data").ITracer() // just make a sublogger

	Tracer.LogEvent("b57e3641").Title("some title 1").Send()

	wg2 := sync.WaitGroup{}
	wg2.Add(2)
	go func() {
		Logger := Tracer.NewSubTracerInitChain("SUBSYSTEM-73853896").ITracer()
		for i := 0; i < 5; i++ {
			Logger.LogEvent("41").Title("g1").Send()
		}
		wg2.Done()
	}()
	go func() {
		Logger := Tracer.NewSubTracerInitChain("5687595c").ITracer()
		for i := 0; i < 5; i++ {
			Logger.LogEvent("48").Title("g2").Send()
		}
		wg2.Done()
	}()
	wg2.Wait()
	Tracer.LogEvent("d0ed6944").Title("complete").Send()
	Tracer.LogEvent("e83ab525").Title("complete").SendMsg("complete")
}

func Test_e1_02(t *testing.T) {

	tracer := loggeterrzl.New(nil).NSTIC("XBegK8r2").ITracer()

	tracer = tracer.NewSubTracerInitChain("SUBSYSTEM-a24487").
		Str("xu", "my-xu").
		Str("ver", "1.2.3").
		ITracer()

	tracer.LogEvent("111b7de0").Title("some title 1").Send()

	var e1 *controlflow.Event

	tracer.SpanBegin("SPAN-3a7ce735")
	{
		tracer.LogEvent("32").SendMsg("hello")
		tracer.SpanBegin("00")
		{
			e1 = tracer.NewControlflowEvent(trace.BacktracePointf("777", "the cause"))

			fmt.Println("--- ", e1)

			tracer.EventEvent("85", e1).Send()
		}
		tracer.SpanEnd()
	}
	tracer.SpanEnd()

	e1.Wrap(trace.BacktracePointf("88", "error"))

	tracer.LogEvent("89").AddControlflowEvent("e1", e1).Send()
}

func Test_e1_03(t *testing.T) {

	tracer := loggeterrzl.New(nil).NSTIC("XBegK8r2").ITracer()

	tracer = tracer.NewSubTracerInitChain("SUBSYSTEM-a24487").
		Str("xu", "my-xu").
		Str("ver", "1.2.3").
		ITracer()

	tracer.LogEvent("111b7de0").Title("some title 1").Send()

	var e1 *controlflow.Event

	tracer.SpanBegin("SPAN-3a7ce735")
	{
		tracer.LogEvent("32").SendMsg("hello")
		tracer.SpanBegin("00")
		{
			err1 := errors.New("some cause error")
			err2 := fmt.Errorf("wrapped; %w", err1)

			e1 = tracer.NewControlflowEvent(&trace.BacktracePoint{"777", err2})

			fmt.Println("--- ", e1)

			tracer.EventEvent("85", e1).Send()
		}
		tracer.SpanEnd()
	}
	tracer.SpanEnd()

	e1.Wrap(trace.BacktracePointf("88", "error"))

	tracer.LogEvent("89").AddControlflowEvent("e1", e1).Send()
}
