package configvalidate

import (
	datavm "github.com/rusriver/nutz/data-vm"
)

type Command [2]uint8

var (
	Command_IsType             = Command{1, 0}
	Command_Exists             = Command{2, 0}
	Command_Absent             = Command{3, 0}
	Command_EqualsEitherString = Command{4, 0}
	Command_EqualsEitherNumber = Command{5, 0}
	Command_NumberIsInRange    = Command{6, 0}
)

const (
	Type_String uint16 = iota
	Type_Integer
	Type_Float
	Type_Bool
	Type_Duration
)

type Instruction struct {
	Id      string
	Command Command // program -> instruction -> command
	Path    []string
	Type    uint16
	Values  []any
}

func (i *Instruction) Execute(s any) (err error) {

	switch i.Command {

	default:
		err = datavm.StdErr(i.Id, i.Path, "wrong command")
		return
	}
	// return
}
