package datavm

import (
	"fmt"
	"strings"
)

type IInstruction interface {
	Execute(s any) error
}

func ExecuteUntilFirstError(s any, program []IInstruction) (err error) {
	for _, i := range program {
		err = i.Execute(s)
		if err != nil {
			return
		}
	}
	return
}

func ExecuteAll(s any, program []IInstruction) (errs []error) {
	for _, i := range program {
		err := i.Execute(s)
		if err != nil {
			errs = append(errs, err)
		}
	}
	return
}

func StdErr(id string, path []string, msg string) (err error) {
	err = fmt.Errorf("instr %v, path %v, error: %v", id, strings.Join(path, "."), msg)
	return
}
