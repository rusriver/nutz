package replicount_test

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/rusriver/nutz/replicount"
	"github.com/rusriver/ttlcache/v3"
)

func Test_001(t *testing.T) {
	cache := ttlcache.New[string, bool]()
	go cache.Start()

	cache.SetWithTTL("a", true, time.Second*2)
	cache.SetWithTTL("s", true, time.Second*4)
	cache.SetWithTTL("d", true, time.Second*6)
	cache.SetWithTTL("f", true, time.Second*8)

	for i := 20; i >= 0; i-- {
		ln := cache.Len()
		keys := cache.Keys()

		fmt.Printf("%v, %v\n", ln, keys)

		time.Sleep(time.Millisecond * 500)
	}
}

func Test_002(t *testing.T) {
	replSet1 := []string{
		"0001",
		"0002",
		"0003",
		"0004a",
	}

	replSet2 := []string{
		"0011",
		"0012",
		"0013",
		"0014",
		"0015",
		"0016",
	}

	t0 := time.Now()

	replicount.New(func(r *replicount.Replicount) {
		r.SlowPollPeriodPerReplica = time.Second * 3
		r.TTLMultiple = 6
		r.FastModeSpeedMultiple = 10

		r.PollFunc = func() (idPtr *string) {
			var id string
			if time.Now().Before(t0.Add(time.Second * 40)) {
				id = replSet1[rand.Intn(len(replSet1))]
			} else {
				id = replSet2[rand.Intn(len(replSet2))]
			}
			fmt.Println("poll", id)
			return &id
		}

		r.ReportResultsFunc = func(newResult *replicount.ChangeableObject) {
			fmt.Printf("result %v %+v\n", newResult.NumberOfReplicas, newResult.ListOfReplicas)
		}

		r.LogFunc = func(s string) {
			fmt.Println(s)
		}
	})

	time.Sleep(time.Second * 90)
}

func Test_003(t *testing.T) {
	replSet1 := []string{
		"0001",
		"0002",
		"0003",
		"0004",
	}

	replSet2 := []string{
		"0011",
		"0012",
		"0013",
		"0014",
		"0015",
		"0016",
	}

	i := 0
	t0 := time.Now()

	replicount.New(func(r *replicount.Replicount) {
		r.SlowPollPeriodPerReplica = time.Second * 3
		r.TTLMultiple = 5
		r.FastModeSpeedMultiple = 10

		r.PollFunc = func() (idPtr *string) {
			var id string
			if time.Now().Before(t0.Add(time.Second * 40)) {
				id = replSet1[i]
				i++
				if i >= len(replSet1) {
					i = 0
				}
			} else {
				id = replSet2[i]
				i++
				if i >= len(replSet2) {
					i = 0
				}
			}
			fmt.Println("poll", id)
			return &id
		}

		r.ReportResultsFunc = func(newResult *replicount.ChangeableObject) {
			fmt.Printf("result %v %+v\n", newResult.NumberOfReplicas, newResult.ListOfReplicas)
		}

		r.LogFunc = func(s string) {
			fmt.Println(s)
		}
	})

	time.Sleep(time.Second * 90)
}

func Test_004(t *testing.T) {
	replSet1 := []string{
		"0001",
		"0002",
		"0003",
		"0004",
	}

	replSet2 := []string{
		"0011",
		"0012",
		"0013",
		"0014",
		"0015",
		"0016",
	}

	i := 0
	t0 := time.Now()

	replicount.New(func(r *replicount.Replicount) {
		r.SlowPollPeriodPerReplica = time.Second * 3
		r.TTLMultiple = 5
		r.FastModeSpeedMultiple = 10

		r.PollFunc = func() (idPtr *string) {
			var id string
			if time.Now().Before(t0.Add(time.Second * 40)) {
				id = replSet1[i]
				i++
				if i >= len(replSet1) {
					i = 0
				}
			} else {
				id = replSet2[i]
				i++
				if i >= len(replSet2) {
					i = 0
				}
			}
			fmt.Println("poll", id)
			time.Sleep(time.Millisecond * 200)
			return &id
		}

		r.ReportResultsFunc = func(newResult *replicount.ChangeableObject) {
			fmt.Printf("result %v %+v\n", newResult.NumberOfReplicas, newResult.ListOfReplicas)
		}

		r.LogFunc = func(s string) {
			fmt.Println(s)
		}
	})

	time.Sleep(time.Second * 120)
}
