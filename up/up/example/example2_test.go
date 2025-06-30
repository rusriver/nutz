package example_test

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand/v2"
	"sync"
	"testing"
	"time"

	"github.com/rusriver/nutz/controlflow"
	loggeterrzl "github.com/rusriver/nutz/loggeterr-zl"
	"github.com/rusriver/nutz/up/fork"
	"github.com/rusriver/nutz/up/reactor"
	"github.com/rusriver/nutz/up/up"
	upcommon "github.com/rusriver/nutz/up/up-common"
	"github.com/rusriver/nutz/up/up/example"
	"github.com/rusriver/nutz/up/watcher"
	"github.com/stretchr/testify/require"
)

// Here we will create a topology with reactors, and will try to use the RecursiveWindUp().
// Consider this a more of a real-life example.

func Test_DI_Container_Example2_00(t *testing.T) {
	DI_Container_Example2_Main(t, 0)
}

func Test_DI_Container_Example2_01(t *testing.T) {
	DI_Container_Example2_Main(t, 1)
}

func DI_Container_Example2_Main(t *testing.T, mode int) {

	/*
		We will model this (the pseudocode) (see also the diagram picture):

			A1 otdg is {
				a1_nguloo := a1_nguloo()
				a1_fakint := a1_fakint()
				a1_mestic := a1_mestic(a1_nguloo, a1_fakint, sparip, robles, histil, petorn, seigul)
			} started by R1 reactor

			R1 depends on subsystems {
				S_myinei == up
				S_nutari == up
				S_uromer == up_or_ceded
			} starts A1

			S_myinei otdg is {
				sostie := new_sostie()
				razinn := new_razinn(sostie)
				seigul := new_seigul(sostie)
				kaffin := new_kaffin(sostie, seigul)
				sparip := new_sparip(sostie, kaffin)
				send R1 the signal we are up
			} started by F1

			S_nutari otdg is {
				robles := new_robles(rebrei, teashi)
				send R1 the signal we are up
			} started by R3 reactor

			S_uromer otdg is {
				petorn := new_petorn()
				histil := new_histil(petorn, itchot, ontras)
				send R1 the signal we are up
			} started by R2 reactor

			R2 depends on subsystems {
				S_obdulu == up
				S_dodism == up
			} starts S_uromer

			R3 depends on subsystems {
				S_armati == up
				S_xylane == up
			} starts S_nutari

			S_obdulu otdg is {
				unborr := new_unborr(aslast, uphili)
				axseps := new_axseps(unborr)
				asdeed := new_asdeed(unborr, axseps)
				itchot := new_itchot(unborr, axseps, asdeed)
				send R2 the signal we are up
			} started by F1

			S_dodism otdg is {
				ontras := new_ontras()
				vexati := new_vexati()
				fopicu := new_fopicu()
				ejacti := new_ejacti()
				hoorgo := new_hoorgo()
				volath := new_volath()
				send R2 the signal we are up
			} started by F1

			S_armati otdg is {
				fatoca := new_fatoca()
				paudat := new_paudat()
				teashi := new_teashi(fatoca, paudat, aslast, uphili)
				send R3 the signal we are up
			} started by F1

			S_xylane otdg is {
				requac := new_requac()
				rebrei := new_rebrei(requac, aslast, uphili)
				send R3 the signal we are up
			} started by F1

			F1 fork is {
				-- start in parallel
				S_obdulu
				S_dodism
				S_armati
				S_xylane
				S_myinei
			} started by S_rastan

			S_rastan otdg is {
				aslast := new_aslast()
				uphili := new_uphili()
				start F1
			} started by ENTRY POINT


		Note that because this model has no diamond topology, it is a DAG, we can safely use
		the RecursiveWindUp() to stop it and remove.
	*/

	/*
		Few words on how to draw it:
			1) draw reactors, forks, otdgs, connect them
			2) draw internals of each otdg in separation
			3) connect cross-otdg dependencies between components
	*/

	/*
		In code, it all starts with creation of reactors (joins) and forks and otdg
		objects, here in main code area (or in infrastructure, if used on a distributed system).

		The same time, define the otdgs constuctors. This is done by defining a subroutines
		(functions). When you will call them, you will actually bring up their contents to life.
		But now you just define those subroutines.

		Note, that otdg is used as a side-object when you explicitly start your components,
		and when you wind up you use the method on the otdg object, it's kind of non-symmetrical,
		worth noting.

		Once you have reactors and forks and otdg objects created and otdg constructors defined,
		you are ready to start it all up. This is done by manually invoking one (or several)
		entry point otdgs. In this example, we have one, the S_rastan.

		The procedure of coding it is somehow mimicking the drawing procedure - first define
		top-level things, create cross-otdg links, the inside of otdgs are coded after that.
	*/

	// Create infra prerequisite components
	tracer := loggeterrzl.New(nil)
	ctx := context.Background()
	wg1 := sync.WaitGroup{}
	wg1.Add(1)

	// define the variables for cross-otdgs.
	// though these are defined here at the top, you typically add them at last stages of coding.
	// you don't need everything here, only those that are used in otdgs other than where these are created.
	var (
		aslast *example.Aslast
		uphili *example.Uphili
		itchot *example.Itchot
		ontras *example.Ontras
		petorn *example.System03
		histil *example.System02
		robles *example.System01
		seigul *example.System03
		sparip *example.System02
	)

	// Create reactors and forks and otdg objects, the top-level topology
	// Start from the entry point, and build from that forward
	S_rastan := &up.OTDG{
		Name:   "S_rastan",
		Tracer: tracer.NSTIC("S_rastan").ITracer(), // must be NSTIC, because may be called in parallel
		Parent: nil,
	}
	F1 := &fork.Fork{
		Name:   "F1",
		Tracer: tracer.NSTIC("F1").ITracer(), // must be NSTIC, because may be called in parallel
		Parent: S_rastan,
	}
	S_myinei := &up.OTDG{
		Name:   "S_myinei",
		Tracer: tracer.NSTIC("S_myinei").ITracer(), // must be NSTIC, because may be called in parallel
		Parent: F1,
	}
	S_armati := &up.OTDG{
		Name:   "S_armati",
		Tracer: tracer.NSTIC("S_armati").ITracer(), // must be NSTIC, because may be called in parallel
		Parent: F1,
	}
	S_xylane := &up.OTDG{
		Name:   "S_xylane",
		Tracer: tracer.NSTIC("S_xylane").ITracer(), // must be NSTIC, because may be called in parallel
		Parent: F1,
	}
	S_obdulu := &up.OTDG{
		Name:   "S_obdulu",
		Tracer: tracer.NSTIC("S_obdulu").ITracer(), // must be NSTIC, because may be called in parallel
		Parent: F1,
	}
	S_dodism := &up.OTDG{
		Name:   "S_dodism",
		Tracer: tracer.NSTIC("S_dodism").ITracer(), // must be NSTIC, because may be called in parallel
		Parent: F1,
	}
	R3 := &reactor.Reactor{
		Name:    "R3",
		Context: ctx,
		Tracer:  tracer.NSTIC("R3").ITracer(), // must be NSTIC, because may be called in parallel
		NeedsPattern: map[string]uint8{
			S_armati.Name: reactor.RequiredStatus_Up,
			S_xylane.Name: reactor.RequiredStatus_Up,
		},
		Func: func(r *reactor.Reactor, msg *reactor.Msg) {

		},
	}
	S_nutari := &up.OTDG{
		Name:   "S_nutari",
		Tracer: tracer.NSTIC("S_nutari").ITracer(), // must be NSTIC, because may be called in parallel
		Parent: R3,
	}
	R2 := &reactor.Reactor{
		Name:    "R2",
		Context: ctx,
		Tracer:  tracer.NSTIC("R2").ITracer(), // must be NSTIC, because may be called in parallel
		NeedsPattern: map[string]uint8{
			S_obdulu.Name: reactor.RequiredStatus_Up,
			S_dodism.Name: reactor.RequiredStatus_Up,
		},
	}
	S_uromer := &up.OTDG{
		Name:   "S_uromer",
		Tracer: tracer.NSTIC("S_uromer").ITracer(), // must be NSTIC, because may be called in parallel
		Parent: R2,
	}
	R1 := &reactor.Reactor{
		Name:    "R1",
		Context: ctx,
		Tracer:  tracer.NSTIC("R1").ITracer(), // must be NSTIC, because may be called in parallel
		NeedsPattern: map[string]uint8{
			S_myinei.Name: reactor.RequiredStatus_Up,
			S_nutari.Name: reactor.RequiredStatus_Up,
			S_uromer.Name: reactor.RequiredStatus_Up,
		},
	}
	R2.Func_CanReact = func(r *reactor.Reactor, msg *reactor.Msg, status *reactor.Status) {
		petorn, histil = Make_S_uromer(S_uromer, ontras, itchot, func(otdg *up.OTDG) {
			R1.SendStatus_Up(otdg, otdg.Name)
		})
	}
	R3.Func_CanReact = func(r *reactor.Reactor, msg *reactor.Msg, status *reactor.Status) {
		Make_S_nutari(S_nutari, func(otdg *up.OTDG) {
			switch mode {
			case 0:
				R1.SendStatus_Up(otdg, otdg.Name)
			case 1:
				R1.SendStatus_Fail(otdg, otdg.Name)
			}
		})
	}
	A1 := &up.OTDG{
		Name:   "A1",
		Tracer: tracer.NSTIC("A1").ITracer(), // must be NSTIC, because may be called in parallel
		Parent: R1,
	}
	R1.Func_CanReact = func(r *reactor.Reactor, msg *reactor.Msg, status *reactor.Status) {
		Make_A1(A1, seigul, sparip, robles, petorn, histil, func(otdg *up.OTDG) {
			wg1.Done()
		})
	}

	// create watchers - you can add them the very last in the coding process
	var finished_by_watchers bool
	W01 := &watcher.Watcher{
		Name:        "W01",
		Tracer:      tracer.NSTIC("W01").ITracer(),
		Context:     ctx,
		Timeout:     time.Second * 15,
		TimeoutFunc: func(w *watcher.Watcher) {},
	}
	W02 := &watcher.Watcher{
		Name:        "W02",
		Tracer:      tracer.NSTIC("W02").ITracer(),
		Context:     ctx,
		Timeout:     time.Second * 15,
		TimeoutFunc: func(w *watcher.Watcher) {},
	}
	W03 := &watcher.Watcher{
		Name:    "W03",
		Tracer:  tracer.NSTIC("W03").ITracer(),
		Context: ctx,
		Timeout: time.Second * 20,
		TimeoutFunc: func(w *watcher.Watcher) {
			A1.RecursiveWindUp(tracer.NSTIC("recursive-wind-up--failure").ITracer())
			finished_by_watchers = true
			wg1.Done()
		},
	}
	W01.Start()
	W02.Start()
	W03.Start()

	// reactors need explicit start, don't forget that
	R1.StartFresh()
	R2.StartFresh()
	R3.StartFresh()

	// the entry point to start the whole thing up
	// make the otdand initiate the whole thing start up process
	aslast, uphili = Make_S_rastan(S_rastan, func(otdg *up.OTDG) {
		F1.RunFuncInParallel(func() {
			seigul, sparip = Make_S_myinei(S_myinei, func(otdg *up.OTDG) {
				R1.SendStatus_Up(otdg, otdg.Name)
			})
		})
		F1.RunFuncInParallel(func() {
			Make_S_armati(S_armati, aslast, uphili, func(otdg *up.OTDG) {
				R3.SendStatus_Up(otdg, otdg.Name)
			})
		})
		F1.RunFuncInParallel(func() {
			Make_S_xylane(S_xylane, aslast, uphili, func(otdg *up.OTDG) {
				R3.SendStatus_Up(otdg, otdg.Name)
			})
		})
		F1.RunFuncInParallel(func() {
			itchot = Make_S_obdulu(S_obdulu, aslast, uphili, func(otdg *up.OTDG) {
				R2.SendStatus_Up(otdg, otdg.Name)
			})
		})
		F1.RunFuncInParallel(func() {
			ontras = Make_S_dodism(S_dodism, func(otdg *up.OTDG) {
				R2.SendStatus_Up(otdg, otdg.Name)
			})
		})
	})

	time.Sleep(time.Second)
	wg1.Wait()
	switch mode {
	case 0:
		require.Equal(t, false, finished_by_watchers)
	case 1:
		require.Equal(t, true, finished_by_watchers)
	}
	if finished_by_watchers {
		tracer.InfoEvent("344").SendMsg("FINISHED BY WATCHERS, TIMED OUT")
		bb, _ := json.MarshalIndent(example.ControlRegistry, "", "    ")
		fmt.Println(string(bb))
		require.Equal(t, true, example.ControlRegistry.MatchThePattern(map[string]int{
			"System01": 0,
			"System02": 0,
			"System03": 0,
			"aslast":   0,
			"itchot":   0,
			"ontras":   0,
			"uphili":   0,
		}))
		fmt.Println("OK, all components are DOWN")
		return
	}

	W01.Stop_Done()
	W02.Stop_Done()
	W03.Stop_Done()

	time.Sleep(time.Second)
	tracer.InfoEvent("287").SendMsg("IS UP")

	// make sure everything is started
	bb, _ := json.MarshalIndent(example.ControlRegistry, "", "    ")
	fmt.Println(string(bb))
	require.Equal(t, true, example.ControlRegistry.MatchThePattern(map[string]int{
		"System01": 12,
		"System02": 7,
		"System03": 5,
		"aslast":   1,
		"itchot":   1,
		"ontras":   1,
		"uphili":   1,
	}))
	var (
		_types int
		_total int
	)
	for _, v := range example.ControlRegistry.Components {
		_types++
		_total += v
	}
	fmt.Printf("OK, %v components of %v types are UP\n", _total, _types)

	time.Sleep(5 * time.Second)

	// we're going to stop then

	// S_rastan.WindUp(tracer.NSTIC("wind-up").ITracer())

	A1.RecursiveWindUp(tracer.NSTIC("recursive-wind-up").ITracer())

	// make sure everything is wind-up
	bb, _ = json.MarshalIndent(example.ControlRegistry, "", "    ")
	fmt.Println(string(bb))
	require.Equal(t, true, example.ControlRegistry.MatchThePattern(map[string]int{
		"System01": 0,
		"System02": 0,
		"System03": 0,
		"aslast":   0,
		"itchot":   0,
		"ontras":   0,
		"uphili":   0,
	}))
	fmt.Println("OK, all components are DOWN")
}

// in real projects, this better be in separate file
// note that all these functions have idiomatic (but not exactly the same) signature
func Make_S_rastan(
	otdg *up.OTDG,
	final_func func(otdg *up.OTDG),
) (
	aslast *example.Aslast,
	uphili *example.Uphili,
) {
	otdg.Tracer.InfoEvent("288").SendMsg("starting")
	controlflow.Try(func() (err any) {
		RandomSleep()

		// syntax example, variant 1
		aslast = up.Up(otdg, "aslast", func(name string) (*example.Aslast, upcommon.IRemovable) {
			s := example.NewAslast(name)
			return s, s
		})

		// syntax example, variant 2
		uphili = otdg.Up("uphili", func(name string) (any, upcommon.IRemovable) {
			s := example.NewUphili(name)
			return s, s
		}).(*example.Uphili)

		otdg.Tracer.InfoEvent("389").SendMsg("done")
		SmallSleep()
		if final_func != nil {
			final_func(otdg)
		}
		return nil
	}).Catch(func(e *controlflow.Event) {
		otdg.Tracer.SmartErrEvent(e, "E320").Send()
	})
	return
}

// in real projects, this better be in separate file
// note that all these functions have idiomatic (but not exactly the same) signature
func Make_S_myinei(
	otdg *up.OTDG,
	final_func func(otdg *up.OTDG),
) (
	seigul *example.System03,
	sparip *example.System02,
) {
	otdg.Tracer.InfoEvent("288").SendMsg("starting")
	controlflow.Try(func() (err any) {
		RandomSleep()

		sostie := up.Up(otdg, "sostie", func(name string) (*example.System01, upcommon.IRemovable) {
			s := example.NewSystem01(name)
			return s, s
		})

		razinn := up.Up(otdg, "razinn", func(name string) (*example.System02, upcommon.IRemovable) {
			s := example.NewSystem02(name, sostie)
			return s, s
		})
		_ = razinn

		seigul = up.Up(otdg, "seigul", func(name string) (*example.System03, upcommon.IRemovable) {
			s := example.NewSystem03(name, sostie)
			return s, s
		})

		kaffin := up.Up(otdg, "kaffin", func(name string) (*example.System01, upcommon.IRemovable) {
			s := example.NewSystem01(name, sostie, seigul)
			return s, s
		})

		sparip = up.Up(otdg, "sparip", func(name string) (*example.System02, upcommon.IRemovable) {
			s := example.NewSystem02(name, sostie, kaffin)
			return s, s
		})

		otdg.Tracer.InfoEvent("389").SendMsg("done")
		SmallSleep()
		if final_func != nil {
			final_func(otdg)
		}
		return nil
	}).Catch(func(e *controlflow.Event) {
		otdg.Tracer.SmartErrEvent(e, "E375").Send()
	})
	return
}

// in real projects, this better be in separate file
// note that all these functions have idiomatic (but not exactly the same) signature
func Make_S_armati(
	otdg *up.OTDG,
	aslast *example.Aslast,
	uphili *example.Uphili,
	final_func func(otdg *up.OTDG),
) {
	otdg.Tracer.InfoEvent("288").SendMsg("starting")
	controlflow.Try(func() (err any) {
		RandomSleep()

		fatoca := up.Up(otdg, "fatoca", func(name string) (*example.System01, upcommon.IRemovable) {
			s := example.NewSystem01(name)
			return s, s
		})

		paudat := up.Up(otdg, "paudat", func(name string) (*example.System02, upcommon.IRemovable) {
			s := example.NewSystem02(name)
			return s, s
		})

		teashi := up.Up(otdg, "teashi", func(name string) (*example.System03, upcommon.IRemovable) {
			s := example.NewSystem03(name, aslast, uphili, fatoca, paudat)
			return s, s
		})
		_ = teashi

		otdg.Tracer.InfoEvent("389").SendMsg("done")
		SmallSleep()
		if final_func != nil {
			final_func(otdg)
		}
		return nil
	}).Catch(func(e *controlflow.Event) {
		otdg.Tracer.SmartErrEvent(e, "E399").Send()
	})
}

// in real projects, this better be in separate file
// note that all these functions have idiomatic (but not exactly the same) signature
func Make_S_xylane(
	otdg *up.OTDG,
	aslast *example.Aslast,
	uphili *example.Uphili,
	final_func func(otdg *up.OTDG),
) {
	otdg.Tracer.InfoEvent("288").SendMsg("starting")
	controlflow.Try(func() (err any) {
		RandomSleep()

		requac := up.Up(otdg, "requac", func(name string) (*example.System01, upcommon.IRemovable) {
			s := example.NewSystem01(name)
			return s, s
		})

		rebrei := up.Up(otdg, "rebrei", func(name string) (*example.System02, upcommon.IRemovable) {
			s := example.NewSystem02(name, requac, aslast, uphili)
			return s, s
		})
		_ = rebrei

		otdg.Tracer.InfoEvent("389").SendMsg("done")
		SmallSleep()
		if final_func != nil {
			final_func(otdg)
		}
		return nil
	}).Catch(func(e *controlflow.Event) {
		otdg.Tracer.SmartErrEvent(e, "E422").Send()
	})
}

// in real projects, this better be in separate file
// note that all these functions have idiomatic (but not exactly the same) signature
func Make_S_obdulu(
	otdg *up.OTDG,
	aslast *example.Aslast,
	uphili *example.Uphili,
	final_func func(otdg *up.OTDG),
) (
	itchot *example.Itchot,
) {
	otdg.Tracer.InfoEvent("288").SendMsg("starting")
	controlflow.Try(func() (err any) {
		RandomSleep()

		unborr := up.Up(otdg, "unborr", func(name string) (*example.System01, upcommon.IRemovable) {
			s := example.NewSystem01(name, aslast, uphili)
			return s, s
		})

		axseps := up.Up(otdg, "axseps", func(name string) (*example.System02, upcommon.IRemovable) {
			s := example.NewSystem02(name, unborr)
			return s, s
		})

		asdeed := up.Up(otdg, "asdeed", func(name string) (*example.System03, upcommon.IRemovable) {
			s := example.NewSystem03(name, axseps, unborr)
			return s, s
		})
		_ = asdeed

		itchot = up.Up(otdg, "itchot", func(name string) (*example.Itchot, upcommon.IRemovable) {
			s := example.NewItchot(name, unborr, axseps)
			return s, s
		})

		otdg.Tracer.InfoEvent("389").SendMsg("done")
		SmallSleep()
		if final_func != nil {
			final_func(otdg)
		}
		return nil
	}).Catch(func(e *controlflow.Event) {
		otdg.Tracer.SmartErrEvent(e, "E447").Send()
	})
	return
}

// in real projects, this better be in separate file
// note that all these functions have idiomatic (but not exactly the same) signature
func Make_S_dodism(
	otdg *up.OTDG,
	final_func func(otdg *up.OTDG),
) (
	ontras *example.Ontras,
) {
	otdg.Tracer.InfoEvent("288").SendMsg("starting")
	controlflow.Try(func() (err any) {
		RandomSleep()

		_ = up.Up(otdg, "fopicu", func(name string) (*example.System01, upcommon.IRemovable) {
			s := example.NewSystem01(name)
			return s, s
		})

		_ = up.Up(otdg, "ejacti", func(name string) (*example.System01, upcommon.IRemovable) {
			s := example.NewSystem01(name)
			return s, s
		})

		_ = up.Up(otdg, "hoorgo", func(name string) (*example.System01, upcommon.IRemovable) {
			s := example.NewSystem01(name)
			return s, s
		})

		_ = up.Up(otdg, "volath", func(name string) (*example.System01, upcommon.IRemovable) {
			s := example.NewSystem01(name)
			return s, s
		})

		_ = up.Up(otdg, "vexati", func(name string) (*example.System01, upcommon.IRemovable) {
			s := example.NewSystem01(name)
			return s, s
		})

		ontras = up.Up(otdg, "ontras", func(name string) (*example.Ontras, upcommon.IRemovable) {
			s := example.NewOntras(name)
			return s, s
		})

		otdg.Tracer.InfoEvent("389").SendMsg("done")
		SmallSleep()
		if final_func != nil {
			final_func(otdg)
		}
		return nil
	}).Catch(func(e *controlflow.Event) {
		otdg.Tracer.SmartErrEvent(e, "E471").Send()
	})
	return
}

// in real projects, this better be in separate file
// note that all these functions have idiomatic (but not exactly the same) signature
func Make_S_nutari(
	otdg *up.OTDG,
	final_func func(otdg *up.OTDG),
) (
	robles *example.System01,
) {
	otdg.Tracer.InfoEvent("288").SendMsg("starting")
	controlflow.Try(func() (err any) {
		RandomSleep()

		robles = up.Up(otdg, "robles", func(name string) (*example.System01, upcommon.IRemovable) {
			s := example.NewSystem01(name)
			return s, s
		})

		otdg.Tracer.InfoEvent("389").SendMsg("done")
		SmallSleep()
		if final_func != nil {
			final_func(otdg)
		}
		return nil
	}).Catch(func(e *controlflow.Event) {
		otdg.Tracer.SmartErrEvent(e, "E495").Send()
	})
	return
}

// in real projects, this better be in separate file
// note that all these functions have idiomatic (but not exactly the same) signature
func Make_S_uromer(
	otdg *up.OTDG,
	ontras *example.Ontras,
	itchot *example.Itchot,
	final_func func(otdg *up.OTDG),
) (
	petorn *example.System03,
	histil *example.System02,
) {
	otdg.Tracer.InfoEvent("288").SendMsg("starting")
	controlflow.Try(func() (err any) {
		RandomSleep()

		petorn = up.Up(otdg, "petorn", func(name string) (*example.System03, upcommon.IRemovable) {
			s := example.NewSystem03(name)
			return s, s
		})

		histil = up.Up(otdg, "histil", func(name string) (*example.System02, upcommon.IRemovable) {
			s := example.NewSystem02(name, petorn, itchot, ontras)
			return s, s
		})

		otdg.Tracer.InfoEvent("389").SendMsg("done")
		SmallSleep()
		if final_func != nil {
			final_func(otdg)
		}
		return nil
	}).Catch(func(e *controlflow.Event) {
		otdg.Tracer.SmartErrEvent(e, "E522").Send()
	})
	return
}

// in real projects, this better be in separate file
// note that all these functions have idiomatic (but not exactly the same) signature
func Make_A1(
	otdg *up.OTDG,
	seigul *example.System03,
	sparip *example.System02,
	robles *example.System01,
	petorn *example.System03,
	histil *example.System02,
	final_func func(otdg *up.OTDG),
) {
	otdg.Tracer.InfoEvent("288").SendMsg("starting")
	controlflow.Try(func() (err any) {
		RandomSleep()

		// syntax example, variant 1
		a1_nguloo := up.Up(otdg, "a1_nguloo", func(name string) (*example.System01, upcommon.IRemovable) {
			s := example.NewSystem01(name)
			return s, s
		})

		// syntax example, variant 2
		a1_fakint := otdg.Up("a1_fakint", func(name string) (any, upcommon.IRemovable) {
			s := example.NewSystem02(name)
			return s, s
		}).(*example.System02)

		up.Up(otdg, "a1_mestic", func(name string) (*example.System03, upcommon.IRemovable) {
			s := example.NewSystem03(name, a1_nguloo, a1_fakint, seigul, sparip, robles, petorn, histil)
			return s, s
		})

		otdg.Tracer.InfoEvent("389").SendMsg("done")
		SmallSleep()
		if final_func != nil {
			final_func(otdg)
		}
		return nil
	}).Catch(func(e *controlflow.Event) {
		otdg.Tracer.SmartErrEvent(e, "E566").Send()
	})
}

func RandomSleep() {
	time.Sleep(time.Duration(rand.IntN(1500))*time.Millisecond + 300*time.Millisecond)
}

func SmallSleep() {
	time.Sleep(300 * time.Millisecond)
}
