package delayed_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	loggeterrzl "github.com/rusriver/nutz/loggeterr-zl"
	"github.com/rusriver/nutz/up/delayed"
)

func Test_Delayed_01(t *testing.T) {

	ctx := context.Background()
	itracer := loggeterrzl.New(nil).NSTIC("fb2b9").ITracer()

	d := &delayed.Delayed{
		Context:  ctx,
		Duration: time.Second * 1,
		Tracer:   itracer.SpanBegin("d1"),
		Func: func(d *delayed.Delayed) {
			d.Tracer.LogEvent("38").SendMsg("fired")
		},
		Func_canceled: func(d *delayed.Delayed) {
			d.Tracer.LogEvent("41").SendMsg("canceled")
		},
	}
	delayed.Start(d)

	fmt.Println("STARTED")

	time.Sleep(time.Second * 2)
	fmt.Println("COMPLETE")
}

func Test_Delayed_02(t *testing.T) {

	itracer := loggeterrzl.New(nil).NSTIC("fb2b9").ITracer()

	ctx, canc := context.WithTimeout(context.Background(), time.Second*1)
	defer canc()

	d := &delayed.Delayed{
		Context:  ctx,
		Duration: time.Second * 5,
		Tracer:   itracer.SpanBegin("d1"),
		Func: func(d *delayed.Delayed) {
			d.Tracer.LogEvent("38").SendMsg("fired")
		},
		Func_canceled: func(d *delayed.Delayed) {
			d.Tracer.LogEvent("41").SendMsg("canceled")
		},
	}
	delayed.Start(d)

	fmt.Println("STARTED")

	time.Sleep(time.Second * 2)
	fmt.Println("COMPLETE")
}

func Test_Delayed_03(t *testing.T) {

	itracer := loggeterrzl.New(nil).NSTIC("fb2b9").ITracer()
	ctx, canc := context.WithTimeout(context.Background(), time.Second*5)

	d := &delayed.Delayed{
		Context:  ctx,
		Duration: time.Second * 5,
		Tracer:   itracer.SpanBegin("d1"),
		Func: func(d *delayed.Delayed) {
			d.Tracer.LogEvent("38").SendMsg("fired")
		},
		Func_canceled: func(d *delayed.Delayed) {
			d.Tracer.LogEvent("41").SendMsg("canceled")
		},
	}
	delayed.Start(d)

	fmt.Println("STARTED")

	time.Sleep(time.Second * 1)
	canc()
	time.Sleep(time.Second * 1)
	fmt.Println("COMPLETE")
}
