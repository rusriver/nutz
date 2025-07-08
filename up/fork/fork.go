package fork

import (
	"sync"

	loggeterrapi "github.com/rusriver/nutz/loggeterr-api"
	"github.com/rusriver/nutz/loggeterr-api/trace"
	upcommon "github.com/rusriver/nutz/up/up-common"
)

type Fork struct {
	Name   string // the name is required, since it's used in a user-side tracer when calling RecursiveWindUp()
	Tracer loggeterrapi.ITracer
	Parent upcommon.IRecursiveWindUp
	count  int
	mu     sync.Mutex
}

func NewFork(name string, tracer loggeterrapi.ITracer, parent upcommon.IRecursiveWindUp) *Fork {
	return &Fork{
		Name:   name,
		Tracer: tracer.NSTIC(name).ITracer(),
		Parent: parent,
	}
}

var _ upcommon.IRecursiveWindUp = (*Fork)(nil)

func (f *Fork) GetName() string {
	return f.Name
}

func (fork *Fork) RunFuncInParallel(f func()) {
	if len(fork.Name) == 0 {
		fork.Tracer.RaiseCause(&trace.BacktracePoint{"0197ada91fcb0d42ba7b", "fork name is empty"})
	}
	fork.mu.Lock()
	defer fork.mu.Unlock()
	fork.count++
	go f()
}

func (fork *Fork) RecursiveWindUp(tracer loggeterrapi.ITracer) {
	tracer.InfoEvent("RWC-f33").SendMsg("rec-windup-call")

	fork.mu.Lock()
	defer fork.mu.Unlock()
	fork.count--
	if fork.count == 0 {
		if fork.Parent != nil {
			tracer.InfoEvent("RWB-39").SendMsg("rec-windup-begin")
			defer tracer.InfoEvent("RWE-40").SendMsg("rec-windup-end")
			fork.Parent.RecursiveWindUp(tracer.NSTIC(fork.Parent.GetName()).ITracer())
		}
	}
}
