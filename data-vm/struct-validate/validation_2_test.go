package structvalidate_test

import (
	"fmt"
	"testing"

	datavm "github.com/rusriver/nutz/data-vm"
	structvalidate "github.com/rusriver/nutz/data-vm/struct-validate"
	"github.com/stretchr/testify/require"
)

func Test_02(t *testing.T) {

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
					Command: structvalidate.Command_ArrayContainsEither,
					Path:    []string{"M", "list"},
					Values:  []any{"a", "s"},
				},
			},
			"",
		},
		{"2.0", data,
			[]datavm.IInstruction{
				&structvalidate.Instruction{
					Id:      "1.0",
					Command: structvalidate.Command_ArrayContainsEither,
					Path:    []string{"M", "list"},
					Values:  []any{"q", "s"},
				},
			},
			"instr 1.0, path M.list, error: contains neither",
		},
		{"3.0", data,
			[]datavm.IInstruction{
				&structvalidate.Instruction{
					Id:      "1.0",
					Command: structvalidate.Command_ArrayContainsEither,
					Path:    []string{"M", "list"},
					Values:  []any{1, 2, 3},
				},
			},
			"",
		},
		{"3.1", data,
			[]datavm.IInstruction{
				&structvalidate.Instruction{
					Id:      "1.0",
					Command: structvalidate.Command_ArrayContainsEither,
					Path:    []string{"M", "list"},
					Values:  []any{1},
				},
			},
			"",
		},
		{"3.2", data,
			[]datavm.IInstruction{
				&structvalidate.Instruction{
					Id:      "1.0",
					Command: structvalidate.Command_ArrayContainsEither,
					Path:    []string{"M", "list"},
					Values:  []any{"1"},
				},
			},
			"",
		},
		{"3.3", data,
			[]datavm.IInstruction{
				&structvalidate.Instruction{
					Id:      "1.0",
					Command: structvalidate.Command_ArrayContainsEither,
					Path:    []string{"M", "list"},
					Values:  []any{},
				},
			},
			"instr 1.0, path M.list, error: contains neither",
		},
		{"3.4", data,
			[]datavm.IInstruction{
				&structvalidate.Instruction{
					Id:      "1.0",
					Command: structvalidate.Command_ArrayContainsEither,
					Path:    []string{"M", "list"},
					Values:  []any{},
				},
			},
			"instr 1.0, path M.list, error: contains neither",
		},
		{"4.1", data,
			[]datavm.IInstruction{
				&structvalidate.Instruction{
					Id:      "1.0",
					Command: structvalidate.Command_EqualsEitherValue,
					Path:    []string{"V5"},
					Values:  []any{false},
				},
			},
			"",
		},
		{"4.2", data,
			[]datavm.IInstruction{
				&structvalidate.Instruction{
					Id:      "1.0",
					Command: structvalidate.Command_EqualsEitherValue,
					Path:    []string{"V5"},
					Values:  []any{"0.1"},
				},
			},
			"instr 1.0, path V5, error: equals to neither",
		},
		{"4.3", data,
			[]datavm.IInstruction{
				&structvalidate.Instruction{
					Id:      "1.0",
					Command: structvalidate.Command_EqualsEitherValue,
					Path:    []string{"V5"},
					Values:  []any{"0"},
				},
			},
			"",
		},
		{"4.4", data,
			[]datavm.IInstruction{
				&structvalidate.Instruction{
					Id:      "1.0",
					Command: structvalidate.Command_EqualsEitherValue,
					Path:    []string{"BoolFalseAny"},
					Values:  []any{"0"},
				},
			},
			"",
		},
		{"5.1", data,
			[]datavm.IInstruction{
				&structvalidate.Instruction{
					Id:      "1.0",
					Command: structvalidate.Command_ArrayContainsAll,
					Path:    []string{"M", "list-2"},
					Values:  []any{"a", "b"},
				},
			},
			"",
		},
		{"5.2", data,
			[]datavm.IInstruction{
				&structvalidate.Instruction{
					Id:      "1.0",
					Command: structvalidate.Command_ArrayContainsAll,
					Path:    []string{"M", "list-2", ""},
					Values:  []any{"a", "b"},
				},
			},
			"instr 1.0, path M.list-2., error: invalid index usage",
		},
	}
	for _, cas := range cases {
		fmt.Println(cas.N)
		err := datavm.ExecuteUntilFirstError(cas.Data, cas.Program, nil)
		errMsg := ""
		if err != nil {
			errMsg = err.Error()
		}
		require.Equal(t, cas.ExpectedMessage, errMsg)
	}
}
