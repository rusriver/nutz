package delayed_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/lithammer/shortuuid/v4"
	"github.com/rusriver/nutz/logger"
	"github.com/rusriver/nutz/up/delayed"
)

func Test_Delayed_01(t *testing.T) {

	ctx := context.Background()

	d := &delayed.Delayed{
		Name:     "d1",
		Context:  ctx,
		Duration: time.Second * 1,
		F: func(d *delayed.Delayed) {
			d.Logger.InfoEvent("38").SendMsg("fired")
		},
		F_canceled: func(d *delayed.Delayed) {
			d.Logger.InfoEvent("41").SendMsg("canceled")
		},
	}
	d.Logger = logger.New(shortuuid.New()[:8], nil).NSLIC("fb2b9").Str("name", d.Name).ILogger()
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
			d.Logger.InfoEvent("38").SendMsg("fired")
		},
		F_canceled: func(d *delayed.Delayed) {
			d.Logger.InfoEvent("41").SendMsg("canceled")
		},
	}
	d.Logger = logger.New(shortuuid.New()[:8], nil).NSLIC("fb2b9").Str("name", d.Name).ILogger()
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
			d.Logger.InfoEvent("38").SendMsg("fired")
		},
		F_canceled: func(d *delayed.Delayed) {
			d.Logger.InfoEvent("41").SendMsg("canceled")
		},
	}
	d.Logger = logger.New(shortuuid.New()[:8], nil).NSLIC("fb2b9").Str("name", d.Name).ILogger()
	delayed.Start(d)

	fmt.Println("STARTED")

	time.Sleep(time.Second * 1)
	canc()
	time.Sleep(time.Second * 1)
	fmt.Println("COMPLETE")
}
