package benchmark

import (
	"fmt"
	"math/rand"
	"os"
	"runtime"
	"testing"

	"github.com/rusriver/nutz/highperformance/geteilterzaehler"
	"github.com/shirou/gopsutil/cpu"
)

// go test --bench=. --benchtime=30s

var Parallelism = 1000 // number of goroutines per one CPU;
var EntropyMod = 32    // for best results, must be x8 less than the breadth

func TestMain(m *testing.M) {
	cpuStat, _ := cpu.Info()
	fmt.Println(cpuStat[0].ModelName, cpuStat[0].Cores, "cores")
	cpus := runtime.GOMAXPROCS(0)
	fmt.Println("GOMAXPROCS:", cpus)

	os.Exit(m.Run())
}

/*
	To get actual results, subtract _infra from _load versions.
*/

func Benchmark_gz_1_infra(b *testing.B) {
	counter := geteilterzaehler.NewInt(func(gz *geteilterzaehler.Int) {
		gz.Breadth = 1
	})

	b.SetParallelism(Parallelism)
	fmt.Println(Parallelism*runtime.GOMAXPROCS(0), "Gs", counter.Breadth, "shards")

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		var wrap_i uint16 = uint16(rand.Intn(5000))
		var x uint16
		for pb.Next() {
			entropy := wrap_i % uint16(EntropyMod)
			x += entropy
			/// payload was here
			wrap_i++
		}
	})
}

func Benchmark_gz_1_load(b *testing.B) {
	counter := geteilterzaehler.NewInt(func(gz *geteilterzaehler.Int) {
		gz.Breadth = 1
	})

	b.SetParallelism(Parallelism)
	fmt.Println(Parallelism*runtime.GOMAXPROCS(0), "Gs", counter.Breadth, "shards")

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		var wrap_i uint16 = uint16(rand.Intn(5000))
		var x uint16
		for pb.Next() {
			entropy := wrap_i % uint16(EntropyMod)
			x += entropy
			counter.ApplyValue(entropy, func(s *geteilterzaehler.Scherbe) {
				s.V++
			})
			wrap_i++
		}
	})
}

func Benchmark_gz_32_infra(b *testing.B) {
	counter := geteilterzaehler.NewInt(func(gz *geteilterzaehler.Int) {
		gz.Breadth = 32
	})

	b.SetParallelism(Parallelism)
	fmt.Println(Parallelism*runtime.GOMAXPROCS(0), "Gs", counter.Breadth, "shards")

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		var wrap_i uint16 = uint16(rand.Intn(5000))
		var x uint16
		for pb.Next() {
			entropy := wrap_i % uint16(EntropyMod)
			x += entropy
			/// payload was here
			wrap_i++
		}
	})
}

func Benchmark_gz_32_load(b *testing.B) {
	counter := geteilterzaehler.NewInt(func(gz *geteilterzaehler.Int) {
		gz.Breadth = 32
	})

	b.SetParallelism(Parallelism)
	fmt.Println(Parallelism*runtime.GOMAXPROCS(0), "Gs", counter.Breadth, "shards")

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		var wrap_i uint16 = uint16(rand.Intn(5000))
		var x uint16
		for pb.Next() {
			entropy := wrap_i % uint16(EntropyMod)
			x += entropy
			counter.ApplyValue(entropy, func(s *geteilterzaehler.Scherbe) {
				s.V++
			})
			wrap_i++
		}
	})
}

func Benchmark_gz_256_infra(b *testing.B) {
	counter := geteilterzaehler.NewInt(func(gz *geteilterzaehler.Int) {
		gz.Breadth = 256
	})

	b.SetParallelism(Parallelism)
	fmt.Println(Parallelism*runtime.GOMAXPROCS(0), "Gs", counter.Breadth, "shards")

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		var wrap_i uint16 = uint16(rand.Intn(5000))
		var x uint16
		for pb.Next() {
			entropy := wrap_i % uint16(EntropyMod)
			x += entropy
			/// payload was here
			wrap_i++
		}
	})
}

func Benchmark_gz_256_load(b *testing.B) {
	counter := geteilterzaehler.NewInt(func(gz *geteilterzaehler.Int) {
		gz.Breadth = 256
	})

	b.SetParallelism(Parallelism)
	fmt.Println(Parallelism*runtime.GOMAXPROCS(0), "Gs", counter.Breadth, "shards")

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		var wrap_i uint16 = uint16(rand.Intn(5000))
		var x uint16
		for pb.Next() {
			entropy := wrap_i % uint16(EntropyMod)
			x += entropy
			counter.ApplyValue(entropy, func(s *geteilterzaehler.Scherbe) {
				s.V++
			})
			wrap_i++
		}
	})
}
