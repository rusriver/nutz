package example1_test

import (
	"fmt"
	"io"
	"sync"
	"testing"

	"github.com/lithammer/shortuuid/v4"
	"github.com/rusriver/configdb/v3"
	"github.com/rusriver/nutz/algoutils"
	"github.com/rusriver/nutz/logger"
	"github.com/rusriver/nutz/logger/loggerapi"
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
	nutzLogger := logger.New(shortuuid.New()[:8], nil, func(rctx *loggerapi.RootContext) {
		rctx.EventOnWireHook = func(p []byte) {
			ChLogsAcc <- p
		}
	})

	var Logger loggerapi.ILogger

	Logger = nutzLogger

	Logger.InfoEvent("d3740cd9").Title("some title 1").Send()
	Logger.InfoEvent("05c1e865").Title("some title 1").Send()

	Logger = nutzLogger.SubLoggerInitChain("a24487").
		Str("xu", "my-xu").
		Str("ver", "1.2.3").
		ILogger()

	Logger.InfoEvent("111b7de0").Title("some title 1").Send()
	Logger.InfoEvent("e3466c79").Title("some title 1").Send()

	Logger = Logger.SubLoggerInitChain("59fc7").ILogger() // just make a sublogger

	Logger.InfoEvent("b57e3641").Title("some title 1").Send()

	wg2 := sync.WaitGroup{}
	wg2.Add(2)
	go func() {
		Logger := Logger.SubLoggerInitChain("73853896").ILogger()
		for i := 0; i < 20; i++ {
			Logger.InfoEvent("41").Title("g1").Send()
		}
		wg2.Done()
	}()
	go func() {
		Logger := Logger.SubLoggerInitChain("5687595c").ILogger()
		for i := 0; i < 20; i++ {
			Logger.InfoEvent("48").Title("g2").Send()
		}
		wg2.Done()
	}()
	wg2.Wait()
	Logger.InfoEvent("d0ed6944").Title("complete").Send()
	Logger.InfoEvent("e83ab525").Title("complete").SendMsg("complete")

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
	nutzLogger := logger.New(shortuuid.New()[:8], nil, func(rctx *loggerapi.RootContext) {
		rctx.EventOnWireHook = func(p []byte) {
			ChLogsAcc <- p
		}
	})

	var Logger loggerapi.ILogger

	Logger = nutzLogger

	Logger.InfoEvent("d3740cd9").Title("some title 1").Send()
	Logger.InfoEvent("05c1e865").Title("some title 1").Send()

	Logger = nutzLogger.SubLoggerInitChain("a24487").
		Str("xu", "my-xu").
		Str("ver", "1.2.3").
		ILogger()

	Logger.InfoEvent("111b7de0").Title("some title 1").Send()
	Logger.InfoEvent("e3466c79").Title("some title 1").Send()

	Logger = Logger.SubLoggerInitChain("59fc7").ILogger() // just make a sublogger

	Logger.InfoEvent("b57e3641").Title("some title 1").Send()

	wg2 := sync.WaitGroup{}
	wg2.Add(2)
	go func() {
		Logger := Logger.SubLoggerInitChain("73853896").ILogger()
		for i := 0; i < howManyIterations; i++ {
			Logger.InfoEvent("41").Title("g1").Send()
		}
		wg2.Done()
	}()
	go func() {
		Logger := Logger.SubLoggerInitChain("5687595c").ILogger()
		for i := 0; i < howManyIterations; i++ {
			Logger.InfoEvent("48").Title("g2").Send()
		}
		wg2.Done()
	}()
	wg2.Wait()
	Logger.InfoEvent("d0ed6944").Title("complete").Send()
	Logger.InfoEvent("e83ab525").Title("complete").SendMsg("complete")

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
		if algoutils.HasSubSequence_Fibers(llinedb.P("msgtags").ListString(), []string{"d0ed6944"}) {
			fmt.Println("Gotcha!")
			controlCriteria["d0ed6944"]++
		}
		if algoutils.HasPrefixSequence_Comparables(llinedb.P("etc").ListInt(), []int{3, 3, 2}) {
			controlCriteria["3,3,2"]++
		}
		if algoutils.HasPrefixSequence_Comparables(llinedb.P("etc").ListInt(), []int{3, 3, 3}) {
			controlCriteria["3,3,3"]++
		}
	}

	fmt.Printf("CONTROL CRITERIA: %+v\n", controlCriteria)
	require.Equal(t, howManyIterations, controlCriteria["3,3,2"])
	require.Equal(t, howManyIterations, controlCriteria["3,3,3"])
	require.Equal(t, 1, controlCriteria["d0ed6944"])
}

func Test_e2_03(t *testing.T) {

	// create a normal logger
	nutzLogger := logger.New(shortuuid.New()[:8], nil, func(rctx *loggerapi.RootContext) {})

	var Logger loggerapi.ILogger

	Logger = nutzLogger

	Logger.InfoEvent("d3740cd9").Title("some title 1").Send()
	Logger.InfoEvent("05c1e865").Title("some title 1").Send()
	RootLogger := Logger

	Logger = Logger.SubLoggerInitChain("a24487").
		Str("xu", "my-xu").
		Str("ver", "1.2.3").
		ILogger()

	// create a sub-logger, with disabled and intercepted output

	howManyIterations := 10000

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

	OfflineLogger := Logger.SubLoggerInitChain("06fca0faf7").WithClonedRootContext(func(rctx *loggerapi.RootContext) {
		rctx.EventOnWireHook = func(p []byte) {
			ChLogsAcc <- p
		}
		rctx.MainWriter = io.Discard
	}).ILogger()

	Logger.InfoEvent("111b7de0").Title("some title 1").Send()
	Logger.InfoEvent("e3466c79").Title("some title 1").Send()

	Logger = Logger.SubLoggerInitChain("59fc7").ILogger() // just make a sublogger

	Logger.InfoEvent("b57e3641").Title("some title 1").Send()

	wg2 := sync.WaitGroup{}
	wg2.Add(4)
	go func() {
		Logger := OfflineLogger.SubLoggerInitChain("73853896").ILogger()
		for i := 0; i < howManyIterations; i++ {
			Logger.InfoEvent("290").Title("g1").Send()
		}
		wg2.Done()
	}()
	go func() {
		Logger := OfflineLogger.SubLoggerInitChain("5687595c").ILogger()
		for i := 0; i < howManyIterations; i++ {
			Logger.InfoEvent("297").Title("g2").Send()
		}
		wg2.Done()
	}()
	go func() {
		Logger := OfflineLogger.SubLoggerInitChain("5141c5c5").ILogger()
		for i := 0; i < howManyIterations; i++ {
			Logger.InfoEvent("304").Title("g3").Send()
		}
		wg2.Done()
	}()
	go func() {
		Logger := OfflineLogger.SubLoggerInitChain("a5042f9c").ILogger()
		for i := 0; i < howManyIterations; i++ {
			Logger.InfoEvent("311").Title("g4").Send()
		}
		wg2.Done()
	}()
	wg2.Wait()
	OfflineLogger.InfoEvent("d0ed6944").Title("complete").Send()
	Logger.InfoEvent("e83ab525").Title("complete").SendMsg("complete")

	close(ChLogsAcc)
	wg1.Wait()

	fmt.Println("===========================================")

	controlCriteria := map[string]int{}

	for i, llinedb := range LogsAccumulated {
		_ = i
		// fmt.Println("LOGLINE", i, string(configdb.MustRenderJson(llinedb.Data, 0)))

		if algoutils.HasSubSequence_Fibers(llinedb.P("msgtags").ListString(), []string{"d0ed6944"}) {
			fmt.Println("Gotcha!")
			controlCriteria["d0ed6944"]++
		}
		if algoutils.HasPrefixSequence_Comparables(llinedb.P("etc").ListInt(), []int{3, 1, 1}) {
			controlCriteria["3,1,1"]++
		}
		if algoutils.HasPrefixSequence_Comparables(llinedb.P("etc").ListInt(), []int{3, 1, 2}) {
			controlCriteria["3,1,2"]++
		}
		if algoutils.HasPrefixSequence_Comparables(llinedb.P("etc").ListInt(), []int{3, 1, 3}) {
			controlCriteria["3,1,3"]++
		}
		if algoutils.HasPrefixSequence_Comparables(llinedb.P("etc").ListInt(), []int{3, 1, 4}) {
			controlCriteria["3,1,4"]++
		}
	}

	fmt.Printf("CONTROL CRITERIA: %+v\n", controlCriteria)
	require.Equal(t, howManyIterations, controlCriteria["3,1,1"])
	require.Equal(t, howManyIterations, controlCriteria["3,1,2"])
	require.Equal(t, howManyIterations, controlCriteria["3,1,3"])
	require.Equal(t, howManyIterations, controlCriteria["3,1,4"])
	require.Equal(t, 1, controlCriteria["d0ed6944"])

	RootLogger.InfoEvent("cd9ef761").Title("complete").SendMsg("complete")

}
