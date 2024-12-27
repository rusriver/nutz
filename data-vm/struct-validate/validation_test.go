package structvalidate_test

import (
	"fmt"
	"testing"

	datavm "github.com/rusriver/nutz/data-vm"
	structvalidate "github.com/rusriver/nutz/data-vm/struct-validate"
	"github.com/stretchr/testify/require"
)

func Test_01(t *testing.T) {

	data := GetData_01()

	cases := []struct {
		N               string
		Data            any
		Program         []datavm.IInstruction
		ExpectedMessage string
	}{
		{"1.0", data,
			[]datavm.IInstruction{
				&structvalidate.Instruction{
					Id:      "1.0",
					Command: structvalidate.Command_IsType,
					Path:    []string{"V1"},
					Type:    structvalidate.Type_String,
				},
				&structvalidate.Instruction{
					Id:      "2.0",
					Command: structvalidate.Command_EqualsEitherValue,
					Path:    []string{"V1"},
					Values:  []any{1, 2, "aa", "root"},
				},
			},
			"",
		},
		{"2.0", data,
			[]datavm.IInstruction{
				&structvalidate.Instruction{
					Id:      "1.0",
					Command: structvalidate.Command_IsType,
					Path:    []string{"V1"},
					Type:    structvalidate.Type_String,
				},
				&structvalidate.Instruction{
					Id:      "1.1",
					Command: structvalidate.Command_IsType,
					Path:    []string{"V2"},
					Type:    structvalidate.Type_Integer,
				},
				&structvalidate.Instruction{
					Id:      "2.0",
					Command: structvalidate.Command_EqualsEitherValue,
					Path:    []string{"V1"},
					Values:  []any{"some"},
				},
			},
			"instr 2.0, path V1, error: equals to neither",
		},
		{"3.0", data,
			[]datavm.IInstruction{
				&structvalidate.Instruction{
					Id:      "1.0",
					Command: structvalidate.Command_IsType,
					Path:    []string{"V1"},
					Type:    structvalidate.Type_Float,
				},
			},
			"instr 1.0, path V1, error: type is not float",
		},
		{"4.0", data,
			[]datavm.IInstruction{
				&structvalidate.Instruction{
					Id:      "1.0",
					Command: structvalidate.Command_Absent,
					Path:    []string{"V0"},
				},
				&structvalidate.Instruction{
					Id:      "2.0",
					Command: structvalidate.Command_Absent,
					Path:    []string{"A", "A", "B", "k2", "V3", "A", "V3"},
				},
			},
			"instr 2.0, path A.A.B.k2.V3.A.V3, error: is not absent",
		},
		{"5.0", data,
			[]datavm.IInstruction{
				&structvalidate.Instruction{
					Id:      "1.0",
					Command: structvalidate.Command_Exists,
					Path:    []string{"A", "A", "B", "k2", "V3", "A", "V3"},
				},
				&structvalidate.Instruction{
					Id:      "2.0",
					Command: structvalidate.Command_Exists,
					Path:    []string{"V0"},
				},
			},
			"instr 2.0, path V0, error: not found",
		},
		{"6.0", data,
			[]datavm.IInstruction{
				&structvalidate.Instruction{
					Id:      "1.0",
					Command: structvalidate.Command_EqualsEitherValue,
					Path:    []string{"A", "A", "B", "k2", "V3", "A", "V1"},
					Values:  []any{"root"},
				},
				&structvalidate.Instruction{
					Id:      "2.0",
					Command: structvalidate.Command_EqualsEitherValue,
					Path:    []string{"A", "A", "B", "k2", "V2"},
					Values:  []any{15},
				},
				&structvalidate.Instruction{
					Id:      "3.0",
					Command: structvalidate.Command_EqualsEitherValue,
					Path:    []string{"A", "A", "B", "k2", "V1"},
					Values:  []any{"k2-V1"},
				},
			},
			"",
		},
		{"7.0", data,
			[]datavm.IInstruction{
				&structvalidate.Instruction{
					Id:      "1.0",
					Command: structvalidate.Command_EqualsEitherValue,
					Path:    []string{"A", "A", "B", "k2", "V3", "A", "V1"},
					Values:  []any{"root"},
				},
				&structvalidate.Instruction{
					Id:      "2.0",
					Command: structvalidate.Command_EqualsEitherValue,
					Path:    []string{"A", "A", "B", "k2", "V2"},
					Values:  []any{15.0},
				},
				&structvalidate.Instruction{
					Id:      "2.1",
					Command: structvalidate.Command_EqualsEitherValue,
					Path:    []string{"A", "A", "B", "k2", "V2"},
					Values:  []any{byte(15)},
				},
				&structvalidate.Instruction{
					Id:      "2.2",
					Command: structvalidate.Command_EqualsEitherValue,
					Path:    []string{"A", "A", "B", "k2", "V2"},
					Values:  []any{"15.0"},
				},
				&structvalidate.Instruction{
					Id:      "3.0",
					Command: structvalidate.Command_EqualsEitherValue,
					Path:    []string{"A", "A", "B", "k3", "V3"},
					Values:  []any{"hello 123"},
				},
			},
			"",
		},
		{"8.0", 1,
			[]datavm.IInstruction{
				&structvalidate.Instruction{
					Id:      "1.0",
					Command: structvalidate.Command_Exists,
					Path:    []string{"a"},
				},
			},
			"instr 1.0, path a, error: not found",
		},
		{"9.0", data,
			[]datavm.IInstruction{
				&structvalidate.Instruction{
					Id:      "1.0",
					Command: structvalidate.Command_ArrayContainsAll,
					Path:    []string{"B", "k4", "V3"},
					Values:  []any{"aaa", "sss"},
				},
			},
			"",
		},
		{"10.0", data,
			[]datavm.IInstruction{
				&structvalidate.Instruction{
					Id:      "1.0",
					Command: structvalidate.Command_ArrayContainsAll,
					Path:    []string{"B", "k4", "V3"},
					Values:  []any{"aaa", "sss", "ddd"},
				},
			},
			"instr 1.0, path B.k4.V3, error: doesn't contain all",
		},
		{"11.0", data,
			[]datavm.IInstruction{
				&structvalidate.Instruction{
					Id:      "1.0",
					Command: structvalidate.Command_ArrayContainsEither,
					Path:    []string{"B", "k4", "V3"},
					Values:  []any{"ddd", "eee", "fff", "sss"},
				},
			},
			"",
		},
		{"12.0", data,
			[]datavm.IInstruction{
				&structvalidate.Instruction{
					Id:      "1.0",
					Command: structvalidate.Command_ArrayContainsEither,
					Path:    []string{"B", "k4", "V3"},
					Values:  []any{"ddd", "eee", "fff"},
				},
			},
			"instr 1.0, path B.k4.V3, error: contains neither",
		},
		{"13.0", data,
			[]datavm.IInstruction{
				&structvalidate.Instruction{
					Id:      "1.0",
					Command: structvalidate.Command_EqualsEitherValue,
					Path:    []string{"V4"},
					Values:  []any{"true"},
				},
			},
			"",
		},
		{"14.0", data,
			[]datavm.IInstruction{
				&structvalidate.Instruction{
					Id:      "1.0",
					Command: structvalidate.Command_EqualsEitherValue,
					Path:    []string{"V4"},
					Values:  []any{1},
				},
			},
			"",
		},
		{"15.0", data,
			[]datavm.IInstruction{
				&structvalidate.Instruction{
					Id:      "1.0",
					Command: structvalidate.Command_EqualsEitherValue,
					Path:    []string{"V4"},
					Values:  []any{0},
				},
			},
			"instr 1.0, path V4, error: equals to neither",
		},
	}
	for _, cas := range cases {
		fmt.Println(cas.N)
		err := datavm.ExecuteUntilFirstError(cas.Data, cas.Program)
		errMsg := ""
		if err != nil {
			errMsg = err.Error()
		}
		require.Equal(t, cas.ExpectedMessage, errMsg)
	}
}
