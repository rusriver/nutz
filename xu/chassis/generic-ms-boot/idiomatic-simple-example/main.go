package main

import (
	"github.com/rusriver/nutz/controlflow"
	"github.com/rusriver/nutz/up/up"
	upcommon "github.com/rusriver/nutz/up/up-common"
	genericmsboot "github.com/rusriver/nutz/xu/chassis/generic-ms-boot"
)

var (
	Version string = "0"
	GitHash string
)

func main() {
	local_start_values := map[string]any{
		"tracer.executable-unit-version":     Version,
		"tracer.executable-unit-fingerprint": GitHash,
	}
	// this is special thing, never wound-up, so outside of otdg
	boot := genericmsboot.NewBoot([]func(b *genericmsboot.Boot){
		func(b *genericmsboot.Boot) {
			b.ProductName = "Idiomatic Simple Example by Someone"
			b.Copyright = "(c) Someone, 2025"
			b.PathToTheTracerConfigSection = "tracer"
			b.L1ConfigFileRequired = true
			b.LocalMapOverride = local_start_values
		},
	})
	tracer := boot.RootTracer
	boot.RootTracer.InfoEvent("BEGIN").Send()

	otdg_main := up.NewOTDG("otdg-main", tracer, nil)

	controlflow.Try(func() (err any) {

		core := up.Up(otdg_main, "core", func(name string) (*S_Core, upcommon.IRemovable) {
			s := NewCore(
				name,
				otdg_main.Tracer.NSTIC(name).ITracer(),
				boot.GetGloveAt(otdg_main.Tracer, name),
				func(s *S_Core) {
					s.SomeUserParam01 = "overridden core--084433f7-4260-4ef8-9cee-fd36f5f1de56"
				},
			)
			return s, s
		})

		up.Up(otdg_main, "sys-1", func(name string) (*S_Sys1, upcommon.IRemovable) {
			s := NewSys1(
				name,
				otdg_main.Tracer.NSTIC(name).ITracer(),
				boot.GetGloveAt(otdg_main.Tracer, name),
				[]func(s *S_Sys1){
					func(s *S_Sys1) {
						s.SomeUserParam01 = "overridden sys1--acdf81f5-1942-4358-a677-60d0b8bd688f"
					},
				},
				core,
			)
			return s, s
		})

		up.Up(otdg_main, "sys-2", func(name string) (*S_Sys2, upcommon.IRemovable) {
			s := NewSys2(
				name,
				otdg_main.Tracer.NSTIC(name).ITracer(),
				boot.GetGloveAt(otdg_main.Tracer, name),
				[]func(s *S_Sys2){
					func(s *S_Sys2) {
						s.SomeUserParam01 = "overridden sys2--07fd2176-247f-4b6d-9039-f8a121cd363d"
					},
				},
				core,
			)
			return s, s
		})

		up.Up(otdg_main, "service", func(name string) (*Service, upcommon.IRemovable) {
			s := NewService(
				name,
				otdg_main.Tracer.NSTIC(name).ITracer(),
				boot.GetGloveAt(otdg_main.Tracer, name),
			)
			return s, s
		})

		return
	}).Catch(func(e *controlflow.Event) {
		otdg_main.Tracer.SmartErrEvent(e, "7b31a5794fbe3d1ca078").Send()
		boot.RootContextCancel()
	})

	tracer.InfoEvent("UP").Send()
	boot.BlockAndWait()
	tracer.InfoEvent("GO-DOWN").Send()
	otdg_main.RecursiveWindUp(tracer.NSTIC("otdg-main-wind-up").ITracer())
	tracer.InfoEvent("END").Send()
}
