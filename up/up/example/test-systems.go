package example

import (
	"time"

	loggeterrapi "github.com/rusriver/nutz/loggeterr-api"
	upcommon "github.com/rusriver/nutz/up/up-common"
)

// ============================

type System01 struct {
	Name string
}

var _ upcommon.IRemovable = (*System01)(nil)

func NewSystem01(name string, deps ...any) *System01 {
	ControlRegistry.Components["System01"]++
	return &System01{
		Name: name,
	}
}

func (s *System01) WindUp(tracer loggeterrapi.ITracer) {
	tracer = tracer.NSTIC(s.Name).ITracer()
	SmallSleep()
	tracer.InfoEvent("26").SendMsg("removed")
	ControlRegistry.Components["System01"]--
}

// ============================

type System02 struct {
	Name string
}

var _ upcommon.IRemovable = (*System02)(nil)

func NewSystem02(name string, deps ...any) *System02 {
	ControlRegistry.Components["System02"]++
	return &System02{
		Name: name,
	}
}

func (s *System02) WindUp(tracer loggeterrapi.ITracer) {
	tracer = tracer.NSTIC(s.Name).ITracer()
	SmallSleep()
	tracer.InfoEvent("26").SendMsg("removed")
	ControlRegistry.Components["System02"]--
}

// ============================

type System03 struct {
	Name string
}

var _ upcommon.IRemovable = (*System03)(nil)

func NewSystem03(name string, deps ...any) *System03 {
	ControlRegistry.Components["System03"]++
	return &System03{
		Name: name,
	}
}

func (s *System03) WindUp(tracer loggeterrapi.ITracer) {
	tracer = tracer.NSTIC(s.Name).ITracer()
	SmallSleep()
	tracer.InfoEvent("26").SendMsg("removed")
	ControlRegistry.Components["System03"]--
}

// ============================

type Aslast struct {
	Name string
}

var _ upcommon.IRemovable = (*Aslast)(nil)

func NewAslast(name string) *Aslast {
	ControlRegistry.Components["aslast"]++
	return &Aslast{
		Name: name,
	}
}

func (s *Aslast) WindUp(tracer loggeterrapi.ITracer) {
	tracer = tracer.NSTIC(s.Name).ITracer()
	SmallSleep()
	tracer.InfoEvent("26").SendMsg("removed")
	ControlRegistry.Components["aslast"]--
}

// ============================

type Uphili struct {
	Name string
}

var _ upcommon.IRemovable = (*Uphili)(nil)

func NewUphili(name string) *Uphili {
	ControlRegistry.Components["uphili"]++
	return &Uphili{
		Name: name,
	}
}

func (s *Uphili) WindUp(tracer loggeterrapi.ITracer) {
	tracer = tracer.NSTIC(s.Name).ITracer()
	SmallSleep()
	tracer.InfoEvent("26").SendMsg("removed")
	ControlRegistry.Components["uphili"]--
}

// ============================

type Itchot struct {
	Name string
}

var _ upcommon.IRemovable = (*Itchot)(nil)

func NewItchot(name string, s1 *System01, s2 *System02) *Itchot {
	ControlRegistry.Components["itchot"]++
	return &Itchot{
		Name: name,
	}
}

func (s *Itchot) WindUp(tracer loggeterrapi.ITracer) {
	tracer = tracer.NSTIC(s.Name).ITracer()
	SmallSleep()
	tracer.InfoEvent("26").SendMsg("removed")
	ControlRegistry.Components["itchot"]--
}

// ============================

type Ontras struct {
	Name string
}

var _ upcommon.IRemovable = (*Ontras)(nil)

func NewOntras(name string) *Ontras {
	ControlRegistry.Components["ontras"]++
	return &Ontras{
		Name: name,
	}
}

func (s *Ontras) WindUp(tracer loggeterrapi.ITracer) {
	tracer = tracer.NSTIC(s.Name).ITracer()
	SmallSleep()
	tracer.InfoEvent("26").SendMsg("removed")
	ControlRegistry.Components["ontras"]--
}

// ======================================================================

func SmallSleep() {
	time.Sleep(300 * time.Millisecond)
}
