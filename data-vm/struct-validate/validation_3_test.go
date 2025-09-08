package structvalidate_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/rusriver/nutz/controlflow"
	datavm "github.com/rusriver/nutz/data-vm"
	executionlog "github.com/rusriver/nutz/data-vm/execution-log"
	structvalidate "github.com/rusriver/nutz/data-vm/struct-validate"
	"github.com/stretchr/testify/require"
)

func Test_03(t *testing.T) {

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
					Command: structvalidate.Command_OR,
					Values: []any{
						[]datavm.IInstruction{
							&structvalidate.Instruction{
								Id:      "1.0.1.1",
								Command: structvalidate.Command_ArrayContainsEither,
								Path:    []string{"M", "list"},
								Values:  []any{"a", "s"},
							},
							&structvalidate.Instruction{
								Id:      "1.0.1.2",
								Command: structvalidate.Command_ArrayContainsEither,
								Path:    []string{"M", "list"},
								Values:  []any{"q", "s"},
							},
							&structvalidate.Instruction{
								Id:      "1.0.1.3",
								Command: structvalidate.Command_ArrayContainsEither,
								Path:    []string{"M", "list"},
								Values:  []any{},
							},
						},
						[]datavm.IInstruction{
							&structvalidate.Instruction{
								Id:      "1.0.2.1",
								Command: structvalidate.Command_ArrayContainsEither,
								Path:    []string{"M", "list"},
								Values:  []any{},
							},
							&structvalidate.Instruction{
								Id:      "1.0.2.2",
								Command: structvalidate.Command_ArrayContainsEither,
								Path:    []string{"M", "list"},
								Values:  []any{"1"},
							},
						},
						[]datavm.IInstruction{
							&structvalidate.Instruction{
								Id:      "1.0.3.1",
								Command: structvalidate.Command_ArrayContainsEither,
								Path:    []string{"M", "list"},
								Values:  []any{1, 2},
							},
							&structvalidate.Instruction{
								Id:      "1.0.3.2",
								Command: structvalidate.Command_ArrayContainsEither,
								Path:    []string{"M", "list"},
								Values:  []any{1},
							},
						},
					},
				},
				&structvalidate.Instruction{
					Id:      "2.0",
					Command: structvalidate.Command_ArrayContainsEither,
					Path:    []string{"M", "list"},
					Values:  []any{1, 2, 3},
				},
			},
			"",
		},
		{"2.0", data,
			[]datavm.IInstruction{
				&structvalidate.Instruction{
					Id:      "1.0",
					Command: structvalidate.Command_OR,
					Values: []any{
						[]datavm.IInstruction{
							&structvalidate.Instruction{
								Id:      "1.0.1.1",
								Command: structvalidate.Command_ArrayContainsEither,
								Path:    []string{"M", "list"},
								Values:  []any{"a", "s"},
							},
							&structvalidate.Instruction{
								Id:      "1.0.1.2",
								Command: structvalidate.Command_ArrayContainsEither,
								Path:    []string{"M", "list"},
								Values:  []any{"q", "s"},
							},
							&structvalidate.Instruction{
								Id:      "1.0.1.3",
								Command: structvalidate.Command_ArrayContainsEither,
								Path:    []string{"M", "list"},
								Values:  []any{},
							},
						},
						[]datavm.IInstruction{
							&structvalidate.Instruction{
								Id:      "1.0.2.1",
								Command: structvalidate.Command_ArrayContainsEither,
								Path:    []string{"M", "list"},
								Values:  []any{},
							},
							&structvalidate.Instruction{
								Id:      "1.0.2.2",
								Command: structvalidate.Command_ArrayContainsEither,
								Path:    []string{"M", "list"},
								Values:  []any{"1"},
							},
						},
						[]datavm.IInstruction{
							&structvalidate.Instruction{
								Id:      "1.0.3.1",
								Command: structvalidate.Command_ArrayContainsEither,
								Path:    []string{"M", "list"},
								Values:  []any{1, 2},
							},
							&structvalidate.Instruction{
								Id:      "1.0.3.2",
								Command: structvalidate.Command_ArrayContainsEither,
								Path:    []string{"M", "list--"},
								Values:  []any{1},
							},
						},
					},
				},
				&structvalidate.Instruction{
					Id:      "2.0",
					Command: structvalidate.Command_ArrayContainsEither,
					Path:    []string{"M", "list"},
					Values:  []any{1, 2, 3},
				},
			},
			"instr 1.0.3.2, path M.list--, error: not found",
		},
		{"3.0", data,
			[]datavm.IInstruction{
				&structvalidate.Instruction{
					Id:      "1.0",
					Command: structvalidate.Command_OR,
					Values: []any{
						[]datavm.IInstruction{
							&structvalidate.Instruction{
								Id:      "1.0.1.1",
								Command: structvalidate.Command_ArrayContainsEither,
								Path:    []string{"M", "list"},
								Values:  []any{"a", "s"},
							},
						},
						[]datavm.IInstruction{
							&structvalidate.Instruction{
								Id:      "1.0.2.1",
								Command: structvalidate.Command_ArrayContainsEither,
								Path:    []string{"M", "list"},
								Values:  []any{},
							},
							&structvalidate.Instruction{
								Id:      "1.0.2.2",
								Command: structvalidate.Command_ArrayContainsEither,
								Path:    []string{"M", "list"},
								Values:  []any{"1"},
							},
						},
						[]datavm.IInstruction{
							&structvalidate.Instruction{
								Id:      "1.0.3.1",
								Command: structvalidate.Command_ArrayContainsEither,
								Path:    []string{"M", "list"},
								Values:  []any{1, 2},
							},
							&structvalidate.Instruction{
								Id:      "1.0.3.2",
								Command: structvalidate.Command_ArrayContainsEither,
								Path:    []string{"M", "list--"},
								Values:  []any{1},
							},
						},
					},
				},
				&structvalidate.Instruction{
					Id:      "2.0",
					Command: structvalidate.Command_ArrayContainsEither,
					Path:    []string{"M", "list"},
					Values:  []any{1, 2, 3},
				},
			},
			"",
		},
	}
	for _, cas := range cases {
		fmt.Println(cas.N)
		xlog := &executionlog.ExecutionLog{}
		err := datavm.ExecuteUntilFirstError(cas.Data, cas.Program, xlog)
		errMsg := ""
		if err != nil {
			errMsg = err.Error()
		}
		require.Equal(t, cas.ExpectedMessage, errMsg)
		fmt.Println(string(controlflow.Must(json.MarshalIndent(xlog, "", "    "))))
	}
}
