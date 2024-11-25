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
		"7cfabc92-d49d-4ad8-9c92-86e01b0f205a",
		"01936542-3a8e-7eb5-b0ab-3044784e9e7d",
		"01936542-463d-7404-80c8-83986bf9df3b",
		"01936542-4eb6-7d80-a37d-c16f4368d4ec",
	}

	for _, s := range ss {
		s2 := []byte(s)
		keyhash.ShuffleBytes(s2, 2)

		fmt.Println(s, "->", string(s2))
	}
}

func Test_ShuffleBytesX_02(t *testing.T) {
	k1 := []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20}
	var k2 []byte

	fmt.Println("---3")
	k2 = make([]byte, len(k1))
	copy(k2, k1)
	keyhash.ShuffleBytesX(k2, 1)
	fmt.Println(len(k1), len(k2))
	fmt.Println(k2)

	fmt.Println("---4")
	k2 = make([]byte, len(k1))
	copy(k2, k1)
	keyhash.ShuffleBytesX(k2, 5)
	fmt.Println(len(k1), len(k2))
	fmt.Println(k2)

	fmt.Println("---5")
	k2 = make([]byte, len(k1))
	copy(k2, k1)
	keyhash.ShuffleBytesX(k2, 15)
	fmt.Println(len(k1), len(k2))
	fmt.Println(k2)

	fmt.Println("---6")
	k2 = make([]byte, len(k1))
	copy(k2, k1)
	keyhash.ShuffleBytesX(k2, 1151)
	fmt.Println(len(k1), len(k2))
	fmt.Println(k2)

	fmt.Println("---7")
	k2 = make([]byte, len(k1))
	copy(k2, k1)
	keyhash.ShuffleBytesX(k2, 1152)
	fmt.Println(len(k1), len(k2))
	fmt.Println(k2)

}

func Test_X_023(t *testing.T) {
	m := map[string]*int{}

	v := 10
	m["a"] = &v

	fmt.Println(m["a"])
	fmt.Println(m["b"])
}
