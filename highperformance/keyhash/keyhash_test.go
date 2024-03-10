package keyhash

import (
	"fmt"
	"math/rand"
	"testing"
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
		hash := Get_8bits([]byte(k), 5)
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
		hash := Get_16bits([]byte(k), 6)
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
		hash := Get_16bits(k, 6)
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
		hash := Get_16bits(k, 6)
		fmt.Printf("%04X\n", hash)
	}
}

func Test_16bits_3(t *testing.T) {
	coverageMap := map[int]struct{}{}

	key := []byte{0, 0, 0, 0, 0, 0}

	for i := 0; i < 10_000_000; i++ {
		for ik := range key {
			key[ik] = byte(rand.Intn(0xFF))
		}
		hash := Get_16bits(key, 6)
		coverageMap[hash] = struct{}{}

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

	fmt.Println(len(coverageMap))
}
