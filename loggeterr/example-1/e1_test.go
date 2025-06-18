package example1_test

import (
	"sync"
	"testing"

	"github.com/lithammer/shortuuid/v4"
	"github.com/rusriver/nutz/loggeterr"
	"github.com/rusriver/nutz/loggeterr/loggeterrapi"
)

func Test_e1_01(t *testing.T) {

	tracer := loggeterr.New(nil).NSTIC(shortuuid.New()[:8]).ITracer()

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
