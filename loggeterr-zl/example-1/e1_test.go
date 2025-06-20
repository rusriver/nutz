package example1_test

import (
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

	Tracer.InfoEvent("d3740cd9").Title("some title 1").Send()
	Tracer.InfoEvent("05c1e865").Title("some title 1").Send()

	Tracer = tracer.NewSubTracerInitChain("SUBSYSTEM-a24487").
		Str("xu", "my-xu").
		Str("ver", "1.2.3").
		ITracer()

	Tracer.InfoEvent("111b7de0").Title("some title 1").Send()
	Tracer.InfoEvent("e3466c79").Title("some title 1").Send()

	Tracer.SpanBegin("SPAN-3a7ce735")
	Tracer.InfoEvent("32").SendMsg("hello")
	Tracer.InfoEvent("33").SendMsg("again")
	Tracer.SpanEnd()

	Tracer.InfoEvent("66afccbd").Title("some title 1").Send()

	Tracer = Tracer.NewSubTracerInitChain("59fc7").Str("some", "data").ITracer() // just make a sublogger

	Tracer.InfoEvent("b57e3641").Title("some title 1").Send()

	wg2 := sync.WaitGroup{}
	wg2.Add(2)
	go func() {
		Logger := Tracer.NewSubTracerInitChain("SUBSYSTEM-73853896").ITracer()
		for i := 0; i < 5; i++ {
			Logger.InfoEvent("41").Title("g1").Send()
		}
		wg2.Done()
	}()
	go func() {
		Logger := Tracer.NewSubTracerInitChain("5687595c").ITracer()
		for i := 0; i < 5; i++ {
			Logger.InfoEvent("48").Title("g2").Send()
		}
		wg2.Done()
	}()
	wg2.Wait()
	Tracer.InfoEvent("d0ed6944").Title("complete").Send()
	Tracer.InfoEvent("e83ab525").Title("complete").SendMsg("complete")
}

func Test_e1_02(t *testing.T) {

	tracer := loggeterrzl.New(nil).NSTIC("XBegK8r2").ITracer()

	tracer = tracer.NewSubTracerInitChain("SUBSYSTEM-a24487").
		Str("xu", "my-xu").
		Str("ver", "1.2.3").
		ITracer()

	tracer.InfoEvent("111b7de0").Title("some title 1").Send()

	var e1 *controlflow.Event

	tracer.SpanBegin("SPAN-3a7ce735")
	{
		tracer.InfoEvent("32").SendMsg("hello")
		tracer.SpanBegin("00")
		{
			e1 = tracer.NewControlflowEvent(trace.BacktracePointf("777", "the cause"))

			tracer.SmartErrEvent(e1, "85").Send()
		}
		tracer.SpanEnd()
	}
	tracer.SpanEnd()

	e1.Wrap(trace.BacktracePointf("88", "error"))

	tracer.InfoEvent("89").AddControlflowEvent("e1", e1).Send()
}
