This is a simple and battle-tested readiness sensor.

Open-sourced now in 2024. I'm not gonna document it, so simple it is.
Just look at the code, its design is obvious.

This sensor can also be used as an inventory runtime registry, or as a liveness or health
sensor for subsystems at runtime.

	-- A small clarification:
		
		In the context of this thing, the "sensor" is the autonomous goroutine (aka "service")
		which accounts for all state and state changes, the "probe" is the piece of code in
		instrumented subsystem, which reports to this sensor. For example, if your application has
		subsystems A, S, D, then you can create this single sensor, then add probes to each of
		A, S, D, and then the sensor will report when all three are ready, or if all three are
		still alive and operational.

		What can be the result of using this sensor? It can trigger an activation of HTTP responder,
		e.g. for Kubernetes, or gRPC responder, it can report prometheus metric, and of course it
		can be used to report the inventory of internal subsystems, et al.

		The "readiness" is when all subsystems of the application are operational, and the application
		is ready to handle the load.

		The "liveness" is periodic check for readiness. Also called the "health-check". Technically
		it's the same thing, only the use differs: when you are interested in readiness, you check
		the app when it starts, and then forget; liveness/health is when you check the app periodically
		all the time.

		Internally, if you want to account for readiness only, you can omit setting the TTL.
		Otherwise, setting the TTL will lead to eviction of subsystem from the list, if it fails
		to report itself periodically.
	--

The DQF is publicly accessible, so you can access underlying registry cache
directly if you want, safely.

--

TODO: add ability to set individual TTL per reported thing, and or values of them.

--

Usage idiom 1:

	sensor, _ = readinesssensor.New(func(c *readinesssensor.Sensor) {
		c.Context = context.Background()
	})

	go func() {
		var state int
		for {
			time.Sleep(time.Second)
			switch state {
			case 0:
				ok, match, missing := sensor.MatchTheProfile([]string{
					"A",
					"B",
					"C",
				})
				if ok {
					logger.InfoEvent().Title("readiness sensor - ALL SYSTEMS READY").
                        Strs("match", match).Send()
					state = 1
				} else {
					logger.InfoEvent().Title("readiness sensor - NOT READY, WAITING").
                        Strs("match", match).Strs("missing", missing).Send()
				}
			case 1:
				metrics.UpAndRunning.Inc()
			}
		}
	}()

	go func() {
		time.Sleep(10 * time.Second)
		sensor.Report("A")

		time.Sleep(10 * time.Second)
		sensor.Report("B")

		time.Sleep(10 * time.Second)
		sensor.Report("C")
	}()

Usage idiom 2:

	sensor, _ = readinesssensor.New(func(c *readinesssensor.Sensor) {
		c.Context = context.Background()
	})

	go func() {
		var state int
		for {
			time.Sleep(time.Second)
			switch state {
			case 0:
				ok, match, missing := sensor.MatchTheProfile([]string{
					"A",
					"B",
					"C",
				})
				if ok {
					logger.InfoEvent().Title("readiness sensor - ALL SYSTEMS READY").
                        Strs("match", match).Send()
					state = 1
				} else {
					logger.InfoEvent().Title("readiness sensor - NOT READY, WAITING").
                        Strs("match", match).Strs("missing", missing).Send()
				}
			case 1:
				metrics.UpAndRunning.Inc()
			}
		}
	}()

    // start service A
    serviceA := a.NewService(
        // options
        sensor.NewProbe("A")
    )

    // start service B
    serviceB := b.NewService(
        // options
        sensor.NewProbe("B")
    )

    // start service C
    serviceC := c.NewService(
        // options
        sensor.NewProbe("C")
    )

Then in any of those services, do:

    s.ReadinessProbe.Ready()

The benefit of this is that you don't need to hardcode the service name inside the service.
Instead, you inject it from the main program.

