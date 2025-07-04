package example_test

import (
	"fmt"
	"testing"

	"github.com/rusriver/nutz/controlflow"
	"github.com/rusriver/nutz/loggeterr-api/trace"
	loggeterrzl "github.com/rusriver/nutz/loggeterr-zl"
	"github.com/rusriver/nutz/up/up"
	upcommon "github.com/rusriver/nutz/up/up-common"
	"github.com/rusriver/nutz/up/up/example"
)

func Test_0001(t *testing.T) {
	type Nodes map[string]any
	n := &Nodes{} // will it init the map?
	(*n)["a"] = 1
	fmt.Printf("%+v\n", *n)
}

func Test_DI_Container_Example1_01(t *testing.T) {

	tracer := loggeterrzl.New(nil) // this is special thing, never wound-up, so outside of otdg

	otdg1 := up.NewOTDG("otdg1", tracer, nil)

	err := controlflow.Try(func() (err any) {
		// syntax example, variant 1
		q := up.Up(otdg1, "Q", func(name string) (*example.System01, upcommon.IRemovable) {
			s := &example.System01{
				Name: name,
			}
			return s, s
		})

		// syntax example, variant 2
		w := otdg1.Up("W", func(name string) (any, upcommon.IRemovable) {
			s := &example.System01{
				Name: name,
			}
			return s, s
		}).(*example.System01)

		e := up.Up(otdg1, "E", func(name string) (*example.System02, upcommon.IRemovable) {
			_ = q
			_ = w
			s := &example.System02{
				Name: name,
			}
			return s, s
		})
		_ = up.Up(otdg1, "R", func(name string) (*example.System03, upcommon.IRemovable) {
			_ = q
			_ = w
			_ = e
			s := &example.System03{
				Name: name,
			}
			return s, s
		})
		return nil
	}).Catch(func(e *controlflow.Event) {
		// -- if e.g. the e failed, then stop q and w
	})
	if err != nil {
		return
	}
	fmt.Println("all up")

	// on stop
	otdg1.WindUp(tracer)
}

func Test_DI_Container_Example1_02(t *testing.T) {

	tracer := loggeterrzl.New(nil) // this is special thing, never wound-up, so outside of otdg

	otdg1 := up.NewOTDG("otdg1", tracer, nil)

	err := controlflow.Try(func() (err any) {
		// syntax example, variant 1
		q := up.Up(otdg1, "Q", func(name string) (*example.System01, upcommon.IRemovable) {
			s := example.NewSystem01(name, otdg1.Tracer.NSTIC(name).ITracer())
			return s, s
		})

		// syntax example, variant 2
		w := otdg1.Up("W", func(name string) (any, upcommon.IRemovable) {
			s := example.NewSystem01(name, otdg1.Tracer.NSTIC(name).ITracer())
			return s, s
		}).(*example.System01)

		e := up.Up(otdg1, "E", func(name string) (*example.System02, upcommon.IRemovable) {
			s := example.NewSystem02(name, otdg1.Tracer.NSTIC(name).ITracer(), q, w)
			if true {
				controlflow.RaiseCause(nil, &trace.BacktracePoint{"504333e6ef10", "E failed"})
			}
			return s, s
		})
		r := up.Up(otdg1, "R", func(name string) (*example.System03, upcommon.IRemovable) {
			s := example.NewSystem03(name, otdg1.Tracer.NSTIC(name).ITracer(), q, w, e)
			return s, s
		})

		_ = r

		return nil
	}).Catch(func(e *controlflow.Event) {
		// -- if e.g. the e failed, then stop q and w
		fmt.Println(e)
		otdg1.WindUp(tracer)
	})
	if err != nil {
		return
	}
	fmt.Println("all up")

	// on stop
	otdg1.WindUp(tracer)
}
