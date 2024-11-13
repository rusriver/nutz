package keyhash_test

import (
	"fmt"
	"testing"

	"github.com/rusriver/nutz/highperformance/keyhash"
)

func Test_ShuffleBytes_1(t *testing.T) {
	ss := []string{
		"aassddffgghh0011223344",
		"1234567890",
	}

	for _, s := range ss {
		s2 := string(keyhash.ShuffleBytes([]byte(s), 5))

		fmt.Println(s, "->", s2)
	}
}

func Test_X_023(t *testing.T) {
	m := map[string]*int{}

	v := 10
	m["a"] = &v

	fmt.Println(m["a"])
	fmt.Println(m["b"])
}
