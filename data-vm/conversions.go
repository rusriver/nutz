package datavm

import (
	"fmt"
	"reflect"
	"strconv"
)

func TypeToUint64(v1 any) (v2 uint64) {
	if v1 == nil {
		v2 = 0
		return
	}
	v1rv := reflect.ValueOf(v1)
	switch v1rv.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		v2 = uint64(v1rv.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		v2 = v1rv.Uint()
	case reflect.Float32, reflect.Float64:
		v2 = uint64(v1rv.Float())
	case reflect.Bool:
		if v1rv.Bool() {
			v2 = 1
		} else {
			v2 = 0
		}
	case reflect.String:
		v3, err := strconv.ParseFloat(v1rv.String(), 64)
		if err == nil {
			v2 = uint64(v3)
		} else {
			v2 = 0
		}
	case reflect.Slice:
		switch v := v1.(type) {
		case []byte:
			v3, err := strconv.ParseFloat(string(v), 64)
			if err == nil {
				v2 = uint64(v3)
			} else {
				v2 = 0
			}
		case []rune:
			v3, err := strconv.ParseFloat(string(v), 64)
			if err == nil {
				v2 = uint64(v3)
			} else {
				v2 = 0
			}
		default:
			v2 = 0
		}
	default:
		v2 = 0
	}
	return
}

func TypeToInt64(v1 any) (v2 int64) {
	if v1 == nil {
		v2 = 0
		return
	}
	v1rv := reflect.ValueOf(v1)
	switch v1rv.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		v2 = v1rv.Int()
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		v2 = int64(v1rv.Uint())
	case reflect.Float32, reflect.Float64:
		v2 = int64(v1rv.Float())
	case reflect.Bool:
		if v1rv.Bool() {
			v2 = 1
		} else {
			v2 = 0
		}
	case reflect.String:
		v3, err := strconv.ParseFloat(v1rv.String(), 64)
		if err == nil {
			v2 = int64(v3)
		} else {
			v2 = 0
		}
	case reflect.Slice:
		switch v := v1.(type) {
		case []byte:
			v3, err := strconv.ParseFloat(string(v), 64)
			if err == nil {
				v2 = int64(v3)
			} else {
				v2 = 0
			}
		case []rune:
			v3, err := strconv.ParseFloat(string(v), 64)
			if err == nil {
				v2 = int64(v3)
			} else {
				v2 = 0
			}
		default:
			v2 = 0
		}
	default:
		v2 = 0
	}
	return
}

func TypeToFloat64(v1 any) (v2 float64) {
	if v1 == nil {
		v2 = 0
		return
	}
	v1rv := reflect.ValueOf(v1)
	switch v1rv.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		v2 = float64(v1rv.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		v2 = float64(v1rv.Uint())
	case reflect.Float32, reflect.Float64:
		v2 = v1rv.Float()
	case reflect.Bool:
		if v1rv.Bool() {
			v2 = 1
		} else {
			v2 = 0
		}
	case reflect.String:
		v3, err := strconv.ParseFloat(v1rv.String(), 64)
		if err == nil {
			v2 = v3
		} else {
			v2 = 0
		}
	case reflect.Slice:
		switch v := v1.(type) {
		case []byte:
			v3, err := strconv.ParseFloat(string(v), 64)
			if err == nil {
				v2 = v3
			} else {
				v2 = 0
			}
		case []rune:
			v3, err := strconv.ParseFloat(string(v), 64)
			if err == nil {
				v2 = v3
			} else {
				v2 = 0
			}
		default:
			v2 = 0
		}
	default:
		v2 = 0
	}
	return
}

func TypeToString(v1 any) (v2 string) {
	if v1 == nil {
		v2 = ""
		return
	}
	v1rv := reflect.ValueOf(v1)
	switch v1rv.Kind() {
	case reflect.String:
		v2 = v1rv.String()
	case reflect.Slice:
		switch v := v1.(type) {
		case []byte:
			v2 = string(v)
		case []rune:
			v2 = string(v)
		default:
			v2 = fmt.Sprint(v1)
		}
	default:
		v2 = fmt.Sprint(v1)
	}
	return
}

func TypeToBool(v1 any) (v2, undef bool) {
	if v1 == nil {
		undef = true
		return
	}
	var err error
	v1rv := reflect.ValueOf(v1)
	switch v1rv.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		v := v1rv.Int()
		v2 = v >= 1
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		v := v1rv.Uint()
		v2 = v >= 1
	case reflect.Float32, reflect.Float64:
		v := v1rv.Float()
		v2 = v >= 1
	case reflect.Bool:
		v2 = v1rv.Bool()
	case reflect.String:
		v2, err = strconv.ParseBool(v1rv.String())
	case reflect.Slice:
		switch v := v1.(type) {
		case []byte:
			v2, err = strconv.ParseBool(string(v))
		case []rune:
			v2, err = strconv.ParseBool(string(v))
		default:
			v2 = false
		}
	default:
		v2 = false
	}
	if err != nil {
		undef = true
	}
	return
}
