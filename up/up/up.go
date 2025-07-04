package up

import (
	"slices"

	loggeterrapi "github.com/rusriver/nutz/loggeterr-api"
	"github.com/rusriver/nutz/loggeterr-api/trace"
	upcommon "github.com/rusriver/nutz/up/up-common"
)

// One-Thread Dependencies Graph
type OTDG struct {
	Name    string // the name is required, since it's used in a user-side tracer when calling RecursiveWindUp()
	Tracer  loggeterrapi.ITracer
	Parent  upcommon.IRecursiveWindUp
	upStack []upcommon.IRemovable
}

func (otdg *OTDG) WindUp(tracer loggeterrapi.ITracer) {
	tracer = tracer.NSTIC(otdg.Name).ITracer()
	tracer.InfoEvent("21").SendMsg("windup-begin")
	defer tracer.InfoEvent("22e").SendMsg("windup-end")
	for _, ir := range slices.Backward(otdg.upStack) {
		ir.WindUp(tracer)
	}
	otdg.upStack = nil
}

// Be careful if your reactor dependency graph has diamond topology (is not a DAG);
// In that case, better implement explicit forward action plan to wind up your system.
func (otdg *OTDG) RecursiveWindUp(tracer loggeterrapi.ITracer) {
	tracer = tracer.NSTIC(otdg.Name).ITracer()
	tracer.InfoEvent("RWB-33").SendMsg("rec-windup-begin")
	defer tracer.InfoEvent("RWE-34").SendMsg("rec-windup-end")
	for _, ir := range slices.Backward(otdg.upStack) {
		ir.WindUp(tracer)
	}
	otdg.upStack = nil
	if otdg.Parent != nil {
		otdg.Parent.RecursiveWindUp(tracer)
	}
}

/*
Usage:

	q := up.Up[*example.System01](otdg1, "Q", func(name string) (*example.System01, up.IRemovable) {
		s := &example.System01{}
		return s, s
	})

	q := up.Up(otdg1, "Q", func(name string) (*example.System01, up.IRemovable) {
		s := &example.System01{}
		return s, s
	})

See also Up() method as equal alternative.
*/
func Up[S any](otdg *OTDG, name string, init_func func(name string) (S, upcommon.IRemovable)) S {
	if len(otdg.Name) == 0 {
		otdg.Tracer.RaiseCause(&trace.BacktracePoint{"019788a6bee773a6-1", "OTDG name is empty"})
	}
	s, ir := init_func(name)
	otdg.upStack = append(otdg.upStack, ir)
	return s
}

/*
Usage:

	w := otdg1.Up("W", func(name string) (any, up.IRemovable) {
		s := &example.System01{}
		return s, s
	}).(*example.System01)

See also Up() function as equal alternative.
*/
func (otdg *OTDG) Up(name string, init_func func(name string) (any, upcommon.IRemovable)) any {
	if len(otdg.Name) == 0 {
		otdg.Tracer.RaiseCause(&trace.BacktracePoint{"019788a6bee773a6-2", "OTDG name is empty"})
	}
	s, ir := init_func(name)
	otdg.upStack = append(otdg.upStack, ir)
	return s
}
