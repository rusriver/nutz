package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/rusriver/nutz/highperformance/geteilterzaehler"
)

func main() {
	ctx, ctxCancel := context.WithCancel(context.Background())

	counter := geteilterzaehler.NewInt(func(gz *geteilterzaehler.Int) {
		gz.Breadth = 32
	})

	var avResult float64

	var av0 float64
	go counter.BackgroundReducer(ctx, time.Millisecond*100, 4, func(s *geteilterzaehler.Scherbe, nextRound bool) {
		if nextRound {
			avResult = av0
			av0 = 0
		}
		av0 += s.V
	})
	go func() {
		tick := time.NewTicker(time.Second)
		for {
			select {
			case <-ctx.Done():
				return
			case <-tick.C:
				fmt.Println("avResult =", avResult)
			}
		}
	}()

	for i := 0; i < 100; i++ {
		go func() {
			var wrap_i uint16 = uint16(rand.Intn(5000))
			for {
				select {
				case <-ctx.Done():
					return
				default:
					time.Sleep(time.Second) // in this case this is better than ticker, because ensures more asynchronicity

					entropy := wrap_i % 0xF
					counter.ApplyValue(entropy, func(s *geteilterzaehler.Scherbe) {
						s.V++
					})
					wrap_i++
				}
			}

		}()
	}

	for i := 0; i < 100; i++ {
		go func() {
			var wrap_i uint16 = uint16(rand.Intn(5000))
			for {
				select {
				case <-ctx.Done():
					return
				default:
					time.Sleep(time.Second) // in this case this is better than ticker, because ensures more asynchronicity

					entropy := wrap_i % 0xF
					counter.ApplyValue(entropy, func(s *geteilterzaehler.Scherbe) {
						s.V--
					})
					wrap_i++
				}
			}

		}()
	}

	time.Sleep(time.Second * 30)
	ctxCancel()
	time.Sleep(time.Second * 1)
	return
}
