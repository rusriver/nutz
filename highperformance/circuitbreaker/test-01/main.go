package main

import (
	"context"
	"log"
	"time"

	"github.com/rusriver/nutz/highperformance/circuitbreaker"
)

func main() {
	cb := circuitbreaker.New(func(cb *circuitbreaker.CircuitBreaker) {
		cb.ChanReportSize = 32
		cb.ChanReigniteSize = 2
		cb.ReignitePeriod = time.Second * 10
		cb.CloseAllowWhenAbove = 3
		cb.OpenDenyWhenBelow = -3
		cb.DenyWaitTime = time.Minute * 1
		cb.Context = context.Background()
		cb.DebugMode = true
	})

	chDQF := make(chan func(), 32)
	tkr01 := time.NewTicker(time.Millisecond * 500)

	G1_Results := true

	go func() {
		for {
			select {
			case f := <-chDQF:
				f()
				for i := 0; i < len(chDQF); i++ {
					f = <-chDQF
					f()
				}

			case <-tkr01.C:
				cb.WantToDo(func() (success bool) {
					log.Println("G1", G1_Results)
					return G1_Results
				})

			} // select
		} // for
	}()

	time.Sleep(time.Second * 10)
	log.Println("+++ 10 sec")
	time.Sleep(time.Second * 5)
	log.Println("+++ 5 sec")

	chDQF <- func() {
		G1_Results = false
		log.Println("+++ G1_Results", G1_Results)
	}
	time.Sleep(time.Second * 30)

	chDQF <- func() {
		G1_Results = true
		log.Println("+++ G1_Results", G1_Results)
	}
	time.Sleep(time.Second * 30)

	chDQF <- func() {
		G1_Results = false
		log.Println("+++ G1_Results", G1_Results)
	}
	time.Sleep(time.Second * 30)

	chDQF <- func() {
		G1_Results = true
		log.Println("+++ G1_Results", G1_Results)
	}
	time.Sleep(time.Second * 90)

}
