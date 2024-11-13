package keyhash_test

import (
	"encoding/binary"
	"fmt"
	"testing"

	"github.com/rusriver/nutz/highperformance/keyhash"
)

func Test_Random_01(t *testing.T) {
	start := 100
	step := 17 // or 13, better prime

	for i := start; i < start+30; i++ {
		k := []byte{0, 0, 0, 0}
		binary.BigEndian.PutUint32(k, uint32(i*step))
		hash := keyhash.UniformDistributionHash16(k)
		fmt.Printf("%04X\n", hash)
	}

	/*
		Pros:
			- very simple determinate pseudo-random gen
			- easy to get different determinate sequences, by simply changing start and step params;

		Cons:
			- must not be used as determinate sequences as really random numbers - here's nothing random;
	*/
}
