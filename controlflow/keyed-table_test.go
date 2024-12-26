package controlflow_test

import (
	"fmt"
	"testing"

	"github.com/rusriver/nutz/controlflow"
)

func Test_02(t *testing.T) {
	sw1 := controlflow.KeyedTable[string]{
		"a": func() int {
			fmt.Println("case a")
			return 0
		},
		"b": func() int {
			fmt.Println("case b")
			return -1
		},
		"c": func() int {
			fmt.Println("case c")
			return 0
		},
		"default": func() int {
			fmt.Println("case default")
			return 0
		},
	}

	fmt.Println("---1")
	sw1.Switch("a", "default")

	fmt.Println("---2")
	sw1.Switch("b", "default")

	fmt.Println("---3")
	sw1.Switch("non-existent", "default")

	fmt.Println("---4")
	sw1.Sequence("a", "b", "default")

	fmt.Println("---5")
	sw1.Sequence("a", "c", "default")

	fmt.Println("---6")
	sw1.Sequence("c", "a", "b")
}
