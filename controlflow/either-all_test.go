package controlflow_test

import (
	"fmt"
	"testing"

	"github.com/rusriver/nutz/controlflow"
)

func Test_Either_01(t *testing.T) {
	if controlflow.Either(
		func() bool {
			return IsSomething1(1, 2, 3)
		},
		func() bool {
			return IsSomething1(4, 4, 4)
		},
	) {
		fmt.Println("EITHER OK")
	}

	/*
		Syntax sugar version:

			if Either(
				IsSomething1(1, 2, 3),
				IsSomething1(4, 4, 4),
			) {
				fmt.Println("EITHER OK")
			}
	*/
}

func Test_All_01(t *testing.T) {
	if controlflow.All(
		func() bool {
			return IsSomething2("a", "a", "a")
		},
		func() bool {
			return IsSomething1(4, 4, 4)
		},
	) {
		fmt.Println("ALL OK")
	}

	/*
		Syntax sugar version:

			if All(
				IsSomething2("a", "a", "a"),
				IsSomething1(4, 4, 4),
			) {
				fmt.Println("EITHER OK")
			}
	*/
}

func IsSomething1(a, s, d int) bool {
	if a == s && s == d {
		return true
	} else {
		return false
	}
}

func IsSomething2(z, x, c string) bool {
	if z == x && x == c {
		return true
	} else {
		return false
	}
}
