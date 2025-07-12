package example

import (
	"time"

	loggeterrapi "github.com/rusriver/nutz/loggeterr-api"
	upcommon "github.com/rusriver/nutz/up/up-common"
)

// ============================

type System01 struct {
	Name   string
	Tracer loggeterrapi.ITracer
}

var _ upcommon.IRemovable = (*System01)(nil)

func NewSystem01(name string, tracer loggeterrapi.ITracer, deps ...any) *System01 {
	ControlRegistry.Components["System01"]++
	tracer.LogEvent("START").Send()
	return &System01{
		Name: name,
	}
}
func (s *System01) GetName() string {
	return s.Name
}
func (s *System01) WindUp(tracer loggeterrapi.ITracer) {
	SmallSleep()
	tracer.LogEvent("26").SendMsg("removed")
	ControlRegistry.Components["System01"]--
}

// ============================

type System02 struct {
	Name   string
	Tracer loggeterrapi.ITracer
}

var _ upcommon.IRemovable = (*System02)(nil)

func NewSystem02(name string, tracer loggeterrapi.ITracer, deps ...any) *System02 {
	ControlRegistry.Components["System02"]++
	tracer.LogEvent("START").Send()
	return &System02{
		Name: name,
	}
}
func (s *System02) GetName() string {
	return s.Name
}
func (s *System02) WindUp(tracer loggeterrapi.ITracer) {
	SmallSleep()
	tracer.LogEvent("26").SendMsg("removed")
	ControlRegistry.Components["System02"]--
}

// ============================

type System03 struct {
	Name   string
	Tracer loggeterrapi.ITracer
}

var _ upcommon.IRemovable = (*System03)(nil)

func NewSystem03(name string, tracer loggeterrapi.ITracer, deps ...any) *System03 {
	ControlRegistry.Components["System03"]++
	tracer.LogEvent("START").Send()
	return &System03{
		Name: name,
	}
}
func (s *System03) GetName() string {
	return s.Name
}
func (s *System03) WindUp(tracer loggeterrapi.ITracer) {
	SmallSleep()
	tracer.LogEvent("26").SendMsg("removed")
	ControlRegistry.Components["System03"]--
}

// ============================

type Aslast struct {
	Name   string
	Tracer loggeterrapi.ITracer
}

var _ upcommon.IRemovable = (*Aslast)(nil)

func NewAslast(name string, tracer loggeterrapi.ITracer) *Aslast {
	ControlRegistry.Components["aslast"]++
	tracer.LogEvent("START").Send()
	return &Aslast{
		Name: name,
	}
}
func (s *Aslast) GetName() string {
	return s.Name
}
func (s *Aslast) WindUp(tracer loggeterrapi.ITracer) {
	SmallSleep()
	tracer.LogEvent("26").SendMsg("removed")
	ControlRegistry.Components["aslast"]--
}

// ============================

type Uphili struct {
	Name   string
	Tracer loggeterrapi.ITracer
}

var _ upcommon.IRemovable = (*Uphili)(nil)

func NewUphili(name string, tracer loggeterrapi.ITracer) *Uphili {
	ControlRegistry.Components["uphili"]++
	tracer.LogEvent("START").Send()
	return &Uphili{
		Name: name,
	}
}
func (s *Uphili) GetName() string {
	return s.Name
}
func (s *Uphili) WindUp(tracer loggeterrapi.ITracer) {
	SmallSleep()
	tracer.LogEvent("26").SendMsg("removed")
	ControlRegistry.Components["uphili"]--
}

// ============================

type Itchot struct {
	Name   string
	Tracer loggeterrapi.ITracer
}

var _ upcommon.IRemovable = (*Itchot)(nil)

func NewItchot(name string, tracer loggeterrapi.ITracer, s1 *System01, s2 *System02) *Itchot {
	ControlRegistry.Components["itchot"]++
	tracer.LogEvent("START").Send()
	return &Itchot{
		Name: name,
	}
}
func (s *Itchot) GetName() string {
	return s.Name
}
func (s *Itchot) WindUp(tracer loggeterrapi.ITracer) {
	SmallSleep()
	tracer.LogEvent("26").SendMsg("removed")
	ControlRegistry.Components["itchot"]--
}

// ============================

type Ontras struct {
	Name   string
	Tracer loggeterrapi.ITracer
}

var _ upcommon.IRemovable = (*Ontras)(nil)

func NewOntras(name string, tracer loggeterrapi.ITracer) *Ontras {
	ControlRegistry.Components["ontras"]++
	tracer.LogEvent("START").Send()
	return &Ontras{
		Name: name,
	}
}
func (s *Ontras) GetName() string {
	return s.Name
}
func (s *Ontras) WindUp(tracer loggeterrapi.ITracer) {
	SmallSleep()
	tracer.LogEvent("26").SendMsg("removed")
	ControlRegistry.Components["ontras"]--
}

// ======================================================================

func SmallSleep() {
	time.Sleep(300 * time.Millisecond)
}
