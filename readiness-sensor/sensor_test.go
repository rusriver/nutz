package readinesssensor_test

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/rusriver/nutz/logger"
	readinesssensor "github.com/rusriver/nutz/readiness-sensor"
)

func Test_Sensor_0111(t *testing.T) {

	wg1 := sync.WaitGroup{}
	wg1.Add(1)

	sensor, _ := readinesssensor.New(func(c *readinesssensor.Sensor) {
		c.Logger = logger.New(nil)
		c.F = func(s *readinesssensor.Sensor) {
			status := s.StatusInternal()
			s.Logger.InfoEvent("23").
				Strs("fresh", status.Fresh).
				Strs("ready", status.Ready).
				Strs("ceded", status.Ceded).
				SendMsg("signalled")
			if status.AllDone {
				s.Logger.InfoEvent("17").SendMsg("ALL SUBSYSTEMS ARE READY")
				wg1.Done()
			}
		}
	})

	p1 := sensor.NewProbe("S1")
	p2 := sensor.NewProbe("S2")
	p3 := sensor.NewProbe("S3")

	time.Sleep(time.Second * 1)

	p1.Ready()

	time.Sleep(time.Second * 1)

	p2.Cede()

	time.Sleep(time.Second * 1)

	p3.Ready()

	wg1.Wait()
	fmt.Println("COMPLETE")
}
