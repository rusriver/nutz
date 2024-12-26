package controlflow_test

import (
	"fmt"
	"testing"

	"github.com/rusriver/nutz/controlflow"
)

func Test_01(t *testing.T) {

	controlflow.MatchTable{
		{func() int {
			return 0
		}, func() int {
			return 0
		}},
		{func() int {
			return 1
		}, func() int {
			fmt.Println("hello")
			return 0
		}},
		{func() int {
			return -1
		}, func() int {
			fmt.Println("hello 2")
			return 0
		}},
	}.Match()

}
