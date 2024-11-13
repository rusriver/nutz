package keyhash_test

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/rusriver/nutz/highperformance/keyhash"
)

func Test_8bits(t *testing.T) {
	keys := []string{
		"100",
		"10019832760",
		"10019832761",
		"10019832762",
		"10019832763",
		"10019832764",
		"10019832765",
		"10019832766",
		"10019832767",
	}

	for _, k := range keys {
		hash := keyhash.Get_8bits([]byte(k), 5)
		fmt.Printf("%02X\n", hash)
	}
}

func Test_16bits(t *testing.T) {
	keys := []string{
		"1",
		"10",
		"100",
		"1000",
		"10000",
		"100000",
		"1000000",
		"10000000",
		"100000000",
		"1000000000",
		"10019832760",
		"10019832761",
		"10019832762",
		"10019832763",
		"10019832764",
		"10019832765",
		"10019832766",
		"10019832767",
	}

	for _, k := range keys {
		hash := keyhash.Get_16bits([]byte(k), 6)
		fmt.Printf("%04X\n", hash)
	}
}

func Test_16bits_1_1_uuid(t *testing.T) {
	keys := []string{
		"92c9bfc2-5e2a-4ef9-8eae-89c724cf8fcd" + "f68458bd-73c5-4a40-a803-bb7631a6c488",
		"92c9bfc2-5e2a-4ef9-8eae-89c724cf8fc4" + "f68458bd-73c5-4a40-a803-bb7631a6c488",
		"92c9bfc2-5e2a-4ef9-8eae-19c724cf8fcd" + "f68458bd-73c5-4a40-a803-bb7631a6c488",
		"92c9bfc2-5e2a-4ef9-8eae-89c724cf8f3d" + "f68458bd-73c5-4a40-a803-bb7631a6c488",
		"dd69b686-5424-4151-b983-bfb0b4fe9904" + "bc0024f9-3719-4a02-a9c9-7f8dec0f8d17",
		"29a09504-9128-4834-bf36-01ec5213f283" + "89677c34-f1b4-4339-b5b0-302b0c9779ac",
		"e06e94bd-1c89-4284-bb6e-6664d0f93f06" + "ff2a75e6-2494-410f-9c65-6fc337758b21",
	}

	for _, k := range keys {
		hash := keyhash.Get_16bits([]byte(k), 0)
		fmt.Printf("%04X\n", hash)
	}
}

func Test_16bits_2(t *testing.T) {
	keys := [][]byte{
		{0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 1},
		{0, 0, 0, 0, 0, 2},
		{0, 0, 0, 0, 0, 3},
		{0, 0, 0, 0, 0, 4},
		{0, 0, 0, 0, 0, 5},
		{0, 0, 0, 0, 0, 6},
		{0, 0, 0, 0, 0, 7},
		{0, 0, 0, 0, 0, 8},
		{0, 0, 0, 0, 0, 9},
		{1, 0, 0, 0, 0, 0},
		{2, 0, 0, 0, 0, 0},
		{3, 0, 0, 0, 0, 0},
		{4, 0, 0, 0, 0, 0},
		{5, 0, 0, 0, 0, 0},
		{6, 0, 0, 0, 0, 0},
		{7, 0, 0, 0, 0, 0},
		{8, 0, 0, 0, 0, 0},
		{9, 0, 0, 0, 0, 0},
		{10, 0, 0, 0, 0, 0},
		{11, 0, 0, 0, 0, 0},
		{12, 0, 0, 0, 0, 0},
		{13, 0, 0, 0, 0, 0},
		{14, 0, 0, 0, 0, 0},
		{15, 0, 0, 0, 0, 0},
		{16, 0, 0, 0, 0, 0},
	}

	for _, k := range keys {
		hash := keyhash.Get_16bits(k, 6)
		fmt.Printf("%04X\n", hash)
	}
}

func Test_16bits_2_5(t *testing.T) {
	keys := [][]byte{
		{0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 1},
		{1, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 1, 0},
		{0, 1, 0, 0, 0, 0},
		{0, 0, 0, 1, 0, 0},
		{0, 0, 1, 0, 0, 0},
	}

	for _, k := range keys {
		hash := keyhash.Get_16bits(k, 6)
		fmt.Printf("%04X\n", hash)
	}
}

func Test_16bits_3(t *testing.T) {
	coverageMap := map[int]int{}
	min, max := -1, 0

	key := []byte{0, 0, 0, 0, 0, 0}

	for i := 0; i < 10_000_000; i++ {
		for ik := range key {
			key[ik] = byte(rand.Intn(0xFF))
		}
		hash := keyhash.Get_16bits(key, 6)
		coverageMap[hash]++

		if len(coverageMap) == 0xFFFF/10*5 {
			fmt.Printf("50%% at %v iterations\n", i)
		}
		if len(coverageMap) == 0xFFFF/10*8 {
			fmt.Printf("80%% at %v iterations\n", i)
		}
		if len(coverageMap) == 0xFFFF {
			fmt.Printf("100%% at %v iterations\n", i)
			break
		}
	}

	fmt.Println("COVERAGE:", len(coverageMap))

	for _, v := range coverageMap {
		if min == -1 || v < min {
			min = v
		}
		if v > max {
			max = v
		}
	}
	fmt.Println("== MIN", min, ", MAX", max, ", DELTA", max-min)
}
