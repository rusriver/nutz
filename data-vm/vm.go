package datavm

import (
	"fmt"
	"strings"

	executionlog "github.com/rusriver/nutz/data-vm/execution-log"
)

type IInstruction interface {
	Execute(s any, xlog *executionlog.ExecutionLog) error
}

func ExecuteUntilFirstError(s any, program []IInstruction, xlog *executionlog.ExecutionLog) (last_err error) {
	for _, i := range program {
		last_err = i.Execute(s, xlog)
		if last_err != nil {
			return
		}
	}
	return
}

func StdErr(id string, path []string, msg string) (err error) {
	err = fmt.Errorf("instr %v, path %v, error: %v", id, strings.Join(path, "."), msg)
	return
}
