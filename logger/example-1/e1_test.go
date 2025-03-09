package example1_test

import (
	"sync"
	"testing"

	"github.com/lithammer/shortuuid/v4"
	"github.com/rusriver/nutz/logger"
	"github.com/rusriver/nutz/logger/loggerapi"
)

func Test_e1_01(t *testing.T) {
	nutzLogger := logger.New(nil, func(rctx *logger.RootContext) {
		rctx.RootId = shortuuid.New()
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

	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		Logger := Logger.SubLoggerInitChain("73853896").ILogger()
		for i := 0; i < 20; i++ {
			Logger.InfoEvent("41").Title("g1").Send()
		}
		wg.Done()
	}()
	go func() {
		Logger := Logger.SubLoggerInitChain("5687595c").ILogger()
		for i := 0; i < 20; i++ {
			Logger.InfoEvent("48").Title("g2").Send()
		}
		wg.Done()
	}()
	wg.Wait()
	Logger.InfoEvent("d0ed6944").Title("complete").Send()
	Logger.InfoEvent("e83ab525").Title("complete").SendMsg("complete")
}
