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
