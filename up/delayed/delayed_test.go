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

	d := &delayed.Delayed{
		Name:     "d1",
		Context:  ctx,
		Duration: time.Second * 1,
		F: func(d *delayed.Delayed) {
			d.Tracer.InfoEvent("38").SendMsg("fired")
		},
		F_canceled: func(d *delayed.Delayed) {
			d.Tracer.InfoEvent("41").SendMsg("canceled")
		},
	}
	d.Tracer = loggeterrzl.New(nil).NSTIC("fb2b9").Str("name", d.Name).ITracer()
	delayed.Start(d)

	fmt.Println("STARTED")

	time.Sleep(time.Second * 2)
	fmt.Println("COMPLETE")
}

func Test_Delayed_02(t *testing.T) {

	ctx, canc := context.WithTimeout(context.Background(), time.Second*1)
	defer canc()

	d := &delayed.Delayed{
		Name:     "d1",
		Context:  ctx,
		Duration: time.Second * 5,
		F: func(d *delayed.Delayed) {
			d.Tracer.InfoEvent("38").SendMsg("fired")
		},
		F_canceled: func(d *delayed.Delayed) {
			d.Tracer.InfoEvent("41").SendMsg("canceled")
		},
	}
	d.Tracer = loggeterrzl.New(nil).NSTIC("fb2b9").Str("name", d.Name).ITracer()
	delayed.Start(d)

	fmt.Println("STARTED")

	time.Sleep(time.Second * 2)
	fmt.Println("COMPLETE")
}

func Test_Delayed_03(t *testing.T) {

	ctx, canc := context.WithTimeout(context.Background(), time.Second*5)

	d := &delayed.Delayed{
		Name:     "d1",
		Context:  ctx,
		Duration: time.Second * 5,
		F: func(d *delayed.Delayed) {
			d.Tracer.InfoEvent("38").SendMsg("fired")
		},
		F_canceled: func(d *delayed.Delayed) {
			d.Tracer.InfoEvent("41").SendMsg("canceled")
		},
	}
	d.Tracer = loggeterrzl.New(nil).NSTIC("fb2b9").Str("name", d.Name).ITracer()
	delayed.Start(d)

	fmt.Println("STARTED")

	time.Sleep(time.Second * 1)
	canc()
	time.Sleep(time.Second * 1)
	fmt.Println("COMPLETE")
}
