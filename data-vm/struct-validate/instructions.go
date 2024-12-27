package structvalidate

import (
	"math"
	"reflect"

	datavm "github.com/rusriver/nutz/data-vm"
	"github.com/rusriver/nutz/lookup"
	"github.com/stretchr/testify/assert"
)

type Command [2]uint8

var (
	Command_IsType              = Command{1, 0}
	Command_Exists              = Command{2, 0}
	Command_Absent              = Command{3, 0}
	Command_EqualsEitherValue   = Command{4, 0}
	Command_ArrayContainsEither = Command{5, 0}
	Command_ArrayContainsAll    = Command{6, 0}
)

const (
	Type_String uint16 = iota
	Type_Integer
	Type_Float
	Type_Bool
)

type Instruction struct {
	Id             string
	Command        Command // program -> instruction -> command
	Path           []string
	Type           uint16
	Values         []any
	FloatPrecision float64
}

func (instr *Instruction) Execute(s any) (err error) {

	switch instr.Command {

	case Command_IsType:
		v, er := lookup.P(s, instr.Path...)
		if er != nil {
			err = datavm.StdErr(instr.Id, instr.Path, er.Error())
			return
		}
		switch instr.Type {
		case Type_String:
			switch v.Kind() {
			case reflect.String:
			default:
				err = datavm.StdErr(instr.Id, instr.Path, "type is not string")
				return
			}
		case Type_Integer:
			switch v.Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			default:
				err = datavm.StdErr(instr.Id, instr.Path, "type is not integer")
				return
			}
		case Type_Float:
			switch v.Kind() {
			case reflect.Float32, reflect.Float64:
			default:
				err = datavm.StdErr(instr.Id, instr.Path, "type is not float")
				return
			}
		case Type_Bool:
			switch v.Kind() {
			case reflect.Bool:
			default:
				err = datavm.StdErr(instr.Id, instr.Path, "type is not bool")
				return
			}
		}

	case Command_Exists:
		_, er := lookup.P(s, instr.Path...)
		if er != nil {
			err = datavm.StdErr(instr.Id, instr.Path, er.Error())
			return
		}

	case Command_Absent:
		_, er := lookup.P(s, instr.Path...)
		if er == nil {
			err = datavm.StdErr(instr.Id, instr.Path, "is not absent")
			return
		}

	case Command_EqualsEitherValue:
		v1, er := lookup.P(s, instr.Path...)
		if er != nil {
			err = datavm.StdErr(instr.Id, instr.Path, er.Error())
			return
		}
		for _, vE := range instr.Values {
			if instr.areEqual_typeAgnostic(v1, vE) {
				return
			}
		}
		err = datavm.StdErr(instr.Id, instr.Path, "equals to neither")
		return

	case Command_ArrayContainsEither:
		v1, er := lookup.P(s, instr.Path...)
		if er != nil {
			err = datavm.StdErr(instr.Id, instr.Path, er.Error())
			return
		}
		switch v1.Kind() {
		case reflect.Array, reflect.Slice:
			for j := 0; j < v1.Len(); j++ {
				v1 := v1.Index(j)
				for _, vE := range instr.Values {
					if instr.areEqual_typeAgnostic(v1, vE) {
						return
					}
				}
			}
		default:
			err = datavm.StdErr(instr.Id, instr.Path, "node must be an array")
			return
		}
		err = datavm.StdErr(instr.Id, instr.Path, "contains neither")
		return

	case Command_ArrayContainsAll:
		v1, er := lookup.P(s, instr.Path...)
		if er != nil {
			err = datavm.StdErr(instr.Id, instr.Path, er.Error())
			return
		}
		switch v1.Kind() {
		case reflect.Array, reflect.Slice:
			found := map[int]bool{}
			for j := 0; j < v1.Len(); j++ {
				v1 := v1.Index(j)
				for vi, vE := range instr.Values {
					if instr.areEqual_typeAgnostic(v1, vE) {
						found[vi] = true
					}
				}
			}
			if len(found) == len(instr.Values) {
				return
			}
		default:
			err = datavm.StdErr(instr.Id, instr.Path, "node must be an array")
			return
		}
		err = datavm.StdErr(instr.Id, instr.Path, "doesn't contain all")
		return

	default:
		err = datavm.StdErr(instr.Id, instr.Path, "wrong command")
		return
	}
	return
}

func (i *Instruction) areEqual_typeAgnostic(v1 reflect.Value, vE any) bool {
	switch v1.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		v1 := datavm.TypeToInt64(v1.Interface())
		vE := datavm.TypeToInt64(vE)
		if v1 == vE {
			return true
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		v1 := datavm.TypeToUint64(v1.Interface())
		vE := datavm.TypeToUint64(vE)
		if v1 == vE {
			return true
		}
	case reflect.Float32, reflect.Float64:
		v1 := datavm.TypeToFloat64(v1.Interface())
		vE := datavm.TypeToFloat64(vE)
		if math.Abs(v1-vE) <= i.FloatPrecision {
			return true
		}
	}
	switch v := v1.Interface().(type) {
	case []byte, []rune, string:
		v1 := datavm.TypeToString(v1.Interface())
		vE := datavm.TypeToString(vE)
		if v1 == vE {
			return true
		}
	case bool:
		vE := datavm.TypeToBool(vE)
		if v == vE {
			return true
		}
	}
	//lint:ignore S1008 ok
	if assert.ObjectsAreEqual(v1.Interface(), vE) {
		return true
	}
	return false
}
