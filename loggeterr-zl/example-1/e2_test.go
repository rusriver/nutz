package example1_test

import (
	"fmt"
	"io"
	"strings"
	"sync"
	"testing"

	"github.com/rusriver/configdb/v3"
	"github.com/rusriver/nutz/algoutils"
	loggeterrapi "github.com/rusriver/nutz/loggeterr-api"
	"github.com/rusriver/nutz/loggeterr-api/accumapi"
	loggeterrzl "github.com/rusriver/nutz/loggeterr-zl"
	"github.com/stretchr/testify/require"
)

func Test_e2_00_1(t *testing.T) {
	ch := make(chan int, 0xF)
	ch <- 1
	ch <- 2
	ch <- 3
	ch <- 4
	ch <- 5
	ch <- 6
	ch <- 7
	ch <- 8
	ch <- 9
	ch <- 10
	ch <- 11
	ch <- 12
	close(ch)
	for m := range ch {
		fmt.Println(m)
	}
}

func Test_e2_00_2(t *testing.T) {
	chNil := make(chan int, 0xF)
	close(chNil)
	chNil = nil

	ch := make(chan int, 0xF)
	ch <- 1
	ch <- 2
	ch <- 3
	ch <- 4
	ch <- 5
	ch <- 6
	ch <- 7
	ch <- 8
	ch <- 9
	ch <- 10
	ch <- 11
	ch <- 12
	close(ch)
L:
	for {
		select {
		case m, ok := <-ch:
			fmt.Println(m, ok)
			if !ok {
				break L
			}
		case <-chNil:
			// see https://stackoverflow.com/a/13666733/11729048
		}
	}
}

func Test_e2_01(t *testing.T) {
	// this is useful idiom for testing
	var LogsAccumulated []*configdb.Root
	ChLogsAcc := make(chan []byte, 100)
	wg1 := sync.WaitGroup{}
	wg1.Add(1)
	go func() {
		for lline := range ChLogsAcc {
			db := (&configdb.InitContext{}).FromBytes(lline).Load()
			LogsAccumulated = append(LogsAccumulated, db)
		}
		wg1.Done()
	}()

	// create a logger
	itracer := loggeterrzl.New(nil, func(rctx *loggeterrapi.RootContext) {
		rctx.EventOnWireHook = func(e loggeterrapi.IEventInHook, p []byte) {
			ChLogsAcc <- p
		}
	}).NSTIC("DZDpuNvt").ITracer()

	var Tracer loggeterrapi.ITracer

	Tracer = itracer

	Tracer.InfoEvent("d3740cd9").Title("some title 1").Send()
	Tracer.InfoEvent("05c1e865").Title("some title 1").Send()

	Tracer = itracer.NewSubTracerInitChain("a24487").
		Str("xu-ver", "1.2.3").
		ITracer()

	Tracer.InfoEvent("111b7de0").Title("some title 1").Send()
	Tracer.InfoEvent("e3466c79").Title("some title 1").Send()

	Tracer = Tracer.NewSubTracerInitChain("59fc7").ITracer() // just make a sublogger

	Tracer.InfoEvent("b57e3641").Title("some title 1").Send()

	wg2 := sync.WaitGroup{}
	wg2.Add(2)
	go func() {
		Logger := Tracer.NewSubTracerInitChain("73853896").ITracer()
		for i := 0; i < 20; i++ {
			Logger.InfoEvent("41").Title("g1").Send()
		}
		wg2.Done()
	}()
	go func() {
		Logger := Tracer.NewSubTracerInitChain("5687595c").ITracer()
		for i := 0; i < 20; i++ {
			Logger.InfoEvent("48").Title("g2").Send()
		}
		wg2.Done()
	}()
	wg2.Wait()
	Tracer.InfoEvent("d0ed6944").Title("complete").Send()
	Tracer.InfoEvent("e83ab525").Title("complete").SendMsg("complete")

	close(ChLogsAcc)
	wg1.Wait()

	fmt.Println("===========================================")

	for i, llinedb := range LogsAccumulated {
		fmt.Println("-- LOG LINE", i, "\n", string(configdb.MustRenderYaml(llinedb.Data)))
		// if llineParsed["etc"].([]float64) == []float64{3, 3, 1} {
		// need:
		// 1: parse [][]byte to []*config.Config
		// 2: match ETC prefix sequence
		// 3: match msgtags prefix sequence
	}

}

func Test_e2_02(t *testing.T) {
	howManyIterations := 20

	// this is useful idiom for testing
	var LogsAccumulated []*configdb.Root
	ChLogsAcc := make(chan []byte, 100)

	// create a logger
	itracer := loggeterrzl.New(nil, func(rctx *loggeterrapi.RootContext) {
		rctx.EventOnWireHook = func(e loggeterrapi.IEventInHook, p []byte) {
			ChLogsAcc <- p
		}
	}).NSTIC("DZDpuNvt").ITracer()

	wg1 := sync.WaitGroup{}
	wg1.Add(1)
	go func() {
		for lline := range ChLogsAcc {
			db := (&configdb.InitContext{}).WithTracer(itracer).FromBytes(lline).Load().DontPanic()
			LogsAccumulated = append(LogsAccumulated, db)
		}
		wg1.Done()
	}()

	var Tracer loggeterrapi.ITracer

	Tracer = itracer

	Tracer.InfoEvent("d3740cd9").Title("some title 1").Send()
	Tracer.InfoEvent("05c1e865").Title("some title 1").Send()

	Tracer = itracer.NewSubTracerInitChain("a24487").
		Str("xu-ver", "1.2.3").
		ITracer()

	Tracer.InfoEvent("111b7de0").Title("some title 1").Send()
	Tracer.InfoEvent("e3466c79").Title("some title 1").Send()

	Tracer = Tracer.NewSubTracerInitChain("59fc7").ITracer() // just make a sublogger

	Tracer.InfoEvent("b57e3641").Title("some title 1").Send()

	wg2 := sync.WaitGroup{}
	wg2.Add(2)
	go func() {
		Logger := Tracer.NewSubTracerInitChain("73853896").ITracer()
		for i := 0; i < howManyIterations; i++ {
			Logger.InfoEvent("41").Title("g1").Send()
		}
		wg2.Done()
	}()
	go func() {
		Logger := Tracer.NewSubTracerInitChain("5687595c").ITracer()
		for i := 0; i < howManyIterations; i++ {
			Logger.InfoEvent("48").Title("g2").Send()
		}
		wg2.Done()
	}()
	wg2.Wait()
	Tracer.InfoEvent("d0ed6944").Title("complete").Send()
	Tracer.InfoEvent("e83ab525").Title("complete").SendMsg("complete")

	close(ChLogsAcc)
	wg1.Wait()

	fmt.Println("===========================================")

	controlCriteria := map[string]int{}

	for i, llinedb := range LogsAccumulated {
		fmt.Println("LOGLINE", i, string(configdb.MustRenderJson(llinedb.Data, 0)))
		// you can:
		// 	- match ETC prefix sequence
		// 	- match msgtags prefix sequence
		//  etc
		if algoutils.HasSubSequence_Fibers(llinedb.P("tags").ListString(), []string{"d0ed6944"}) {
			fmt.Println("Gotcha!")
			controlCriteria["d0ed6944"]++
		}
		if algoutils.HasPrefixSequence_Comparables(llinedb.P("etc").ListInt(), []int{1, 3, 3, 2}) {
			controlCriteria["1,3,3,2"]++
		}
		if algoutils.HasPrefixSequence_Comparables(llinedb.P("etc").ListInt(), []int{1, 3, 3, 3}) {
			controlCriteria["1,3,3,3"]++
		}
	}

	fmt.Printf("CONTROL CRITERIA: %+v\n", controlCriteria)
	require.Equal(t, howManyIterations, controlCriteria["1,3,3,2"])
	require.Equal(t, howManyIterations, controlCriteria["1,3,3,3"])
	require.Equal(t, 1, controlCriteria["d0ed6944"])
}

func Test_e2_03(t *testing.T) {

	// create a normal logger
	itracer := loggeterrzl.New(nil, func(rctx *loggeterrapi.RootContext) {}).NSTIC("DZDpuNvt").ITracer()

	var Tracer loggeterrapi.ITracer

	Tracer = itracer

	Tracer.InfoEvent("d3740cd9").Title("some title 1").Send()
	Tracer.InfoEvent("05c1e865").Title("some title 1").Send()
	RootLogger := Tracer

	Tracer = Tracer.NewSubTracerInitChain("_TLK").
		Str("xu-ver", "1.2.3").
		ITracer()

	// create a sub-logger, with disabled and intercepted output

	howManyIterations := 10000

	// this is useful idiom for testing
	var LogsAccumulated []*configdb.Root
	ChLogsAcc := make(chan []byte, 100)
	wg1 := sync.WaitGroup{}
	wg1.Add(1)
	go func() {
		for lline := range ChLogsAcc {
			db := (&configdb.InitContext{}).WithTracer(Tracer).FromBytes(lline).Load().DontPanic()
			LogsAccumulated = append(LogsAccumulated, db)
		}
		wg1.Done()
	}()

	OfflineTracer := Tracer.NewSubTracerInitChain("06fca0faf7").WithClonedRootContext(func(rctx *loggeterrapi.RootContext) {
		rctx.EventOnWireHook = func(e loggeterrapi.IEventInHook, p []byte) {
			ChLogsAcc <- p
		}
		rctx.MainWriter = io.Discard
	}).ITracer()

	Tracer.InfoEvent("111b7de0").Title("some title 1").Send()
	Tracer.InfoEvent("e3466c79").Title("some title 1").Send()

	Tracer = Tracer.NewSubTracerInitChain("59fc7").ITracer() // just make a sublogger

	Tracer.InfoEvent("b57e3641").Title("some title 1").Send()

	wg2 := sync.WaitGroup{}
	wg2.Add(4)
	go func() {
		Logger := OfflineTracer.NewSubTracerInitChain("73853896").ITracer()
		for i := 0; i < howManyIterations; i++ {
			Logger.InfoEvent("290").Title("g1").Send()
		}
		wg2.Done()
	}()
	go func() {
		Logger := OfflineTracer.NewSubTracerInitChain("5687595c").ITracer()
		for i := 0; i < howManyIterations; i++ {
			Logger.InfoEvent("297").Title("g2").Send()
		}
		wg2.Done()
	}()
	go func() {
		Logger := OfflineTracer.NewSubTracerInitChain("5141c5c5").ITracer()
		for i := 0; i < howManyIterations; i++ {
			Logger.InfoEvent("304").Title("g3").Send()
		}
		wg2.Done()
	}()
	go func() {
		Logger := OfflineTracer.NewSubTracerInitChain("a5042f9c").ITracer()
		for i := 0; i < howManyIterations; i++ {
			Logger.InfoEvent("311").Title("g4").Send()
		}
		wg2.Done()
	}()
	wg2.Wait()
	OfflineTracer.InfoEvent("d0ed6944").Title("complete").Send()
	Tracer.InfoEvent("e83ab525").Title("complete").SendMsg("complete")

	close(ChLogsAcc)
	wg1.Wait()

	fmt.Println("===========================================")

	controlCriteria := map[string]int{}

	for i, llinedb := range LogsAccumulated {
		_ = i
		// fmt.Println("LOGLINE", i, string(configdb.MustRenderJson(llinedb.Data, 0)))

		if algoutils.HasSubSequence_Fibers(llinedb.P("tags").ListString(), []string{"d0ed6944"}) {
			fmt.Println("Gotcha!")
			controlCriteria["d0ed6944"]++
		}
		if algoutils.HasPrefixSequence_Comparables(llinedb.P("etc").ListInt(), []int{1, 3, 1, 1}) {
			controlCriteria["1,3,1,1"]++
		}
		if algoutils.HasPrefixSequence_Comparables(llinedb.P("etc").ListInt(), []int{1, 3, 1, 2}) {
			controlCriteria["1,3,1,2"]++
		}
		if algoutils.HasPrefixSequence_Comparables(llinedb.P("etc").ListInt(), []int{1, 3, 1, 3}) {
			controlCriteria["1,3,1,3"]++
		}
		if algoutils.HasPrefixSequence_Comparables(llinedb.P("etc").ListInt(), []int{1, 3, 1, 4}) {
			controlCriteria["1,3,1,4"]++
		}
	}

	fmt.Printf("CONTROL CRITERIA: %+v\n", controlCriteria)
	require.Equal(t, howManyIterations, controlCriteria["1,3,1,1"])
	require.Equal(t, howManyIterations, controlCriteria["1,3,1,2"])
	require.Equal(t, howManyIterations, controlCriteria["1,3,1,3"])
	require.Equal(t, howManyIterations, controlCriteria["1,3,1,4"])
	require.Equal(t, 1, controlCriteria["d0ed6944"])

	RootLogger.InfoEvent("cd9ef761").Title("complete").SendMsg("complete")

}

func Test_e2_04(t *testing.T) {

	// create a normal logger
	itracer := loggeterrzl.New(nil, func(rctx *loggeterrapi.RootContext) {}).NSTIC("DZDpuNvt").ITracer()

	var Tracer loggeterrapi.ITracer

	Tracer = itracer

	Tracer.InfoEvent("d3740cd9").Title("some title 1").Send()
	Tracer.InfoEvent("05c1e865").Title("some title 1").Send()
	RootLogger := Tracer

	Tracer = Tracer.NewSubTracerInitChain("_TLK").
		Str("xu-ver", "1.2.3").
		ITracer()

	// create a sub-logger, with disabled and intercepted output

	howManyIterations := 10000

	it1 := Tracer.NewSubTracerInitChain("06fca0faf7").ITracer()
	var OfflineTracer accumapi.ILogsAccumulator = it1.(*loggeterrzl.Tracer) // because of stupid cycle import ban in Go
	OfflineTracer.SwitchOnTheOfflineAccumulation()

	Tracer.InfoEvent("111b7de0").Title("some title 1").Send()
	Tracer.InfoEvent("e3466c79").Title("some title 1").Send()

	Tracer = Tracer.NewSubTracerInitChain("59fc7").ITracer() // just make a sublogger

	Tracer.InfoEvent("b57e3641").Title("some title 1").Send()

	wg2 := sync.WaitGroup{}
	wg2.Add(4)
	go func() {
		Logger := OfflineTracer.NewSubTracerInitChain("73853896").ITracer()
		for i := 0; i < howManyIterations; i++ {
			Logger.InfoEvent("290").Title("g1").Send()
		}
		wg2.Done()
	}()
	go func() {
		Logger := OfflineTracer.NewSubTracerInitChain("5687595c").ITracer()
		for i := 0; i < howManyIterations; i++ {
			Logger.InfoEvent("297").Title("g2").Send()
		}
		wg2.Done()
	}()
	go func() {
		Logger := OfflineTracer.NewSubTracerInitChain("5141c5c5").ITracer()
		for i := 0; i < howManyIterations; i++ {
			Logger.InfoEvent("304").Title("g3").Send()
		}
		wg2.Done()
	}()
	go func() {
		Logger := OfflineTracer.NewSubTracerInitChain("a5042f9c").ITracer()
		for i := 0; i < howManyIterations; i++ {
			Logger.InfoEvent("311").Title("g4").Send()
		}
		wg2.Done()
	}()
	wg2.Wait()
	OfflineTracer.InfoEvent("d0ed6944").Title("complete").Send()
	Tracer.InfoEvent("e83ab525").Title("complete").SendMsg("complete")

	LogsAccumulated := OfflineTracer.FinishAccumulationAndFetch()

	fmt.Println("===========================================")

	controlCriteria := map[string]int{}

	for i, llinedb := range LogsAccumulated {
		_ = i
		// fmt.Println("LOGLINE", i, string(configdb.MustRenderJson(llinedb.Data, 0)))

		if algoutils.HasSubSequence_Fibers(llinedb.P("tags").ListString(), []string{"d0ed6944"}) {
			fmt.Println("Gotcha!")
			controlCriteria["d0ed6944"]++
		}
		if algoutils.HasPrefixSequence_Comparables(llinedb.P("etc").ListInt(), []int{1, 3, 1, 1}) {
			controlCriteria["1,3,1,1"]++
		}
		if algoutils.HasPrefixSequence_Comparables(llinedb.P("etc").ListInt(), []int{1, 3, 1, 2}) {
			controlCriteria["1,3,1,2"]++
		}
		if algoutils.HasPrefixSequence_Comparables(llinedb.P("etc").ListInt(), []int{1, 3, 1, 3}) {
			controlCriteria["1,3,1,3"]++
		}
		if algoutils.HasPrefixSequence_Comparables(llinedb.P("etc").ListInt(), []int{1, 3, 1, 4}) {
			controlCriteria["1,3,1,4"]++
		}
	}

	fmt.Printf("CONTROL CRITERIA: %+v\n", controlCriteria)
	require.Equal(t, howManyIterations, controlCriteria["1,3,1,1"])
	require.Equal(t, howManyIterations, controlCriteria["1,3,1,2"])
	require.Equal(t, howManyIterations, controlCriteria["1,3,1,3"])
	require.Equal(t, howManyIterations, controlCriteria["1,3,1,4"])
	require.Equal(t, 1, controlCriteria["d0ed6944"])

	RootLogger.InfoEvent("cd9ef761").Title("complete").SendMsg("complete")

	fmt.Println("===========================================")

	// now try to reuse the same logger for accumulation again

	OfflineTracer.RestartTheAccumulation()

	for i := 0; i < 5; i++ {
		OfflineTracer.InfoEvent("c67f9c71").Send()
	}

	// you can even use the Zerologger directly, and its output will be accumulated too:
	zerologger := OfflineTracer.(*loggeterrzl.Tracer).ZeroLoggerPtr
	zerologger.Info().Msgf("Hello from Zerologger!")

	LogsAccumulated = OfflineTracer.FinishAccumulationAndFetch()
	for i, llinedb := range LogsAccumulated {
		_ = i
		fmt.Println("LOGLINE", i, string(configdb.MustRenderJson(llinedb.Data, 0)))
	}
}

func Test_e2_05(t *testing.T) {

	// create a normal logger
	itracer := loggeterrzl.New(nil, func(rctx *loggeterrapi.RootContext) {}).NSTIC("DZDpuNvt").ITracer()

	var Tracer loggeterrapi.ITracer

	Tracer = itracer

	Tracer.InfoEvent("d3740cd9").Title("some title 1").Send()
	Tracer.InfoEvent("05c1e865").Title("some title 1").Send()
	RootLogger := Tracer

	Tracer = Tracer.NewSubTracerInitChain("_TLK").
		Str("xu-ver", "1.2.3").
		ITracer()

	// create a sub-logger, with disabled and intercepted output

	howManyIterations := 10000

	ol1 := Tracer.NewSubTracerInitChain("06fca0faf7").ITracer()
	var OfflineTracer accumapi.ILogsAccumulator = ol1.(*loggeterrzl.Tracer) // because of stupid cycle import ban in Go
	OfflineTracer.SwitchOnTheOfflineAccumulation()

	Tracer.InfoEvent("111b7de0").Title("some title 1").Send()
	Tracer.InfoEvent("e3466c79").Title("some title 1").Send()

	Tracer = Tracer.NewSubTracerInitChain("59fc7").ITracer() // just make a sublogger

	Tracer.InfoEvent("b57e3641").Title("some title 1").Send()

	wg2 := sync.WaitGroup{}
	wg2.Add(4)
	go func() {
		Logger := OfflineTracer.NewSubTracerInitChain("73853896").ITracer()
		for i := 0; i < howManyIterations; i++ {
			Logger.InfoEvent("290").Title("g1").Send()
		}
		wg2.Done()
	}()
	go func() {
		Logger := OfflineTracer.NewSubTracerInitChain("5687595c").ITracer()
		for i := 0; i < howManyIterations; i++ {
			Logger.InfoEvent("297").Title("g2").Send()
		}
		wg2.Done()
	}()
	go func() {
		Logger := OfflineTracer.NewSubTracerInitChain("5141c5c5").ITracer()
		for i := 0; i < howManyIterations; i++ {
			Logger.InfoEvent("304").Title("g3").Send()
		}
		wg2.Done()
	}()
	go func() {
		Logger := OfflineTracer.NewSubTracerInitChain("a5042f9c").ITracer()
		for i := 0; i < howManyIterations; i++ {
			Logger.InfoEvent("311").Title("g4").Send()
		}
		wg2.Done()
	}()
	wg2.Wait()
	OfflineTracer.InfoEvent("d0ed6944").Title("complete").Send()
	Tracer.InfoEvent("e83ab525").Title("complete").SendMsg("complete")

	LogsAccumulated := OfflineTracer.FinishAccumulationAndFetch()

	fmt.Println("===========================================")

	controlCriteria := map[string]int{}

	for i, llinedb := range LogsAccumulated {
		_ = i
		// fmt.Println("LOGLINE", i, string(configdb.MustRenderJson(llinedb.Data, 0)))

		// c59145cfdde048f5a3ef

		fmtt_msgtags := algoutils.FMTTF[bool, string]{
			"d0ed6944": {func(kk []string) bool {
				fmt.Println("Gotcha!")
				controlCriteria[strings.Join(kk, " ")]++
				return true
			}, nil},
		}

		fmtt_etc := algoutils.FMTTF[bool, int]{
			3: {nil, algoutils.FMTTF[bool, int]{
				1: {nil, algoutils.FMTTF[bool, int]{
					1: {func(kk []int) bool {
						controlCriteria[fmt.Sprint(kk)]++
						return true
					}, nil},
					2: {func(kk []int) bool {
						controlCriteria[fmt.Sprint(kk)]++
						return true
					}, nil},
					3: {func(kk []int) bool {
						controlCriteria[fmt.Sprint(kk)]++
						return true
					}, nil},
					4: {func(kk []int) bool {
						controlCriteria[fmt.Sprint(kk)]++
						return true
					}, nil},
				}},
			}},
		}

		fmtt_msgtags.FastMatch(llinedb.P("tags").ListString())

		fmtt_etc.FastMatch(llinedb.P("etc").ListInt())
	}

	fmt.Printf("CONTROL CRITERIA: %+v\n", controlCriteria)
	require.Equal(t, howManyIterations, controlCriteria["[3 1 1]"])
	require.Equal(t, howManyIterations, controlCriteria["[3 1 2]"])
	require.Equal(t, howManyIterations, controlCriteria["[3 1 3]"])
	require.Equal(t, howManyIterations, controlCriteria["[3 1 4]"])
	require.Equal(t, 1, controlCriteria["d0ed6944"])

	RootLogger.InfoEvent("cd9ef761").Title("complete").SendMsg("complete")

	fmt.Println("===========================================")

	// now try to reuse the same logger for accumulation again

	OfflineTracer.RestartTheAccumulation()

	for i := 0; i < 5; i++ {
		OfflineTracer.InfoEvent("c67f9c71").Send()
	}

	// you can even use the Zerologger directly, and its output will be accumulated too:
	zerologger := OfflineTracer.(*loggeterrzl.Tracer).ZeroLoggerPtr
	zerologger.Info().Msgf("Hello from Zerologger!")

	LogsAccumulated = OfflineTracer.FinishAccumulationAndFetch()
	for i, llinedb := range LogsAccumulated {
		_ = i
		fmt.Println("LOGLINE", i, string(configdb.MustRenderJson(llinedb.Data, 0)))
	}
}
