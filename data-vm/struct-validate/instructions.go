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
	Command_IsType            = Command{1, 0}
	Command_Exists            = Command{2, 0}
	Command_Absent            = Command{3, 0}
	Command_EqualsEitherValue = Command{4, 0}
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

func (i *Instruction) Execute(s any) (err error) {

	switch i.Command {

	case Command_IsType:
		v, er := lookup.P(s, i.Path...)
		if er != nil {
			err = datavm.StdErr(i.Id, i.Path, er.Error())
			return
		}
		switch i.Type {
		case Type_String:
			switch v.Kind() {
			case reflect.String:
			default:
				err = datavm.StdErr(i.Id, i.Path, "type is not string")
				return
			}
		case Type_Integer:
			switch v.Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			default:
				err = datavm.StdErr(i.Id, i.Path, "type is not integer")
				return
			}
		case Type_Float:
			switch v.Kind() {
			case reflect.Float32, reflect.Float64:
			default:
				err = datavm.StdErr(i.Id, i.Path, "type is not float")
				return
			}
		case Type_Bool:
			switch v.Kind() {
			case reflect.Bool:
			default:
				err = datavm.StdErr(i.Id, i.Path, "type is not bool")
				return
			}
		}

	case Command_Exists:
		_, er := lookup.P(s, i.Path...)
		if er != nil {
			err = datavm.StdErr(i.Id, i.Path, er.Error())
			return
		}

	case Command_Absent:
		_, er := lookup.P(s, i.Path...)
		if er == nil {
			err = datavm.StdErr(i.Id, i.Path, "is not absent")
			return
		}

	case Command_EqualsEitherValue:
		v1, er := lookup.P(s, i.Path...)
		if er != nil {
			err = datavm.StdErr(i.Id, i.Path, er.Error())
			return
		}
		for _, vE := range i.Values {
			switch v1.Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				v1 := datavm.TypeToInt64(v1.Interface())
				vE := datavm.TypeToInt64(vE)
				if v1 == vE {
					return
				}
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				v1 := datavm.TypeToUint64(v1.Interface())
				vE := datavm.TypeToUint64(vE)
				if v1 == vE {
					return
				}
			case reflect.Float32, reflect.Float64:
				v1 := datavm.TypeToFloat64(v1.Interface())
				vE := datavm.TypeToFloat64(vE)
				if math.Abs(v1-vE) <= i.FloatPrecision {
					return
				}
			}
			switch v1.Interface().(type) {
			case []byte, []rune, string:
				v1 := datavm.TypeToString(v1.Interface())
				vE := datavm.TypeToString(vE)
				if v1 == vE {
					return
				}
			}
			if assert.ObjectsAreEqual(v1.Interface(), vE) {
				return
			}
		}
		err = datavm.StdErr(i.Id, i.Path, "equals to neither")
		return

	default:
		err = datavm.StdErr(i.Id, i.Path, "wrong command")
		return
	}
	return
}
