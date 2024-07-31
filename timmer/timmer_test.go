package timmer_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/rusriver/nutz/timmer"
)

func Test_001(t *testing.T) {
	fmt.Println("1")
	t0 := timmer.New()
	t0.Restart(time.Second * 2)
	<-t0.C
	fmt.Println("2")
	t0.Restart(time.Second * 3)
	<-t0.C
	fmt.Println("3")
	t0.Restart(time.Second * 4)
	fmt.Println("4")
	t0.Stop()
	fmt.Println("5")
	t0.Restart(time.Second * 6)
	<-t0.C
	fmt.Println("6")
}
