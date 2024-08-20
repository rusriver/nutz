This is a simple and battle-tested readiness sensor.

Open-sourced now in 2024. I'm not gonna document it, so simple it is.
Just look at the code, its design is obvious.

This sensor can also be used as an inventory runtime registry, or as a liveness or health
sensor for subsystems at runtime.

The DQF is publicly accessible, so you can access underlying registry cache
directly if you want, safely.

--

TODO: add ability to set individual TTL per reported thing, and or values of them.

--

Usage idiom:

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

