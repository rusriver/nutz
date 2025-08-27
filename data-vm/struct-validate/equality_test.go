package structvalidate_test

import (
	"fmt"
	"testing"

	structvalidate "github.com/rusriver/nutz/data-vm/struct-validate"
	"github.com/rusriver/nutz/lookup"
	"github.com/stretchr/testify/require"
)

type X int

const (
	X1 X = 0
	X2 X = 2
	X3 X = 3
)

type BigStruct struct {
	Sub struct {
		X X
	}
}

func Test_Equality_001(t *testing.T) {

	big := &BigStruct{}
	big.Sub.X = X2

	val, err := lookup.DotP(big, "Sub.X")
	fmt.Println("--31", err, val, val.Kind())

	instr := &structvalidate.Instruction{
		FloatPrecision: 0.001,
	}
	var x bool

	x = instr.AreEqual_typeAgnostic(val, 0)
	fmt.Println("--39", x, ", expected false")
	require.Equal(t, x, false)

	x = instr.AreEqual_typeAgnostic(val, 1)
	fmt.Println("--42", x, ", expected false")
	require.Equal(t, x, false)

	x = instr.AreEqual_typeAgnostic(val, 2)
	fmt.Println("--45", x, ", expected true")
	require.Equal(t, x, true)
}
