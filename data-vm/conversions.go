package datavm

import (
	"fmt"
	"strconv"
)

func TypeToUint64(v1 any) (v2 uint64) {
	switch v := v1.(type) {
	case int:
		v2 = uint64(v)
	case int8:
		v2 = uint64(v)
	case int16:
		v2 = uint64(v)
	case int32:
		v2 = uint64(v)
	case int64:
		v2 = uint64(v)
	case uint:
		v2 = uint64(v)
	case uint8:
		v2 = uint64(v)
	case uint16:
		v2 = uint64(v)
	case uint32:
		v2 = uint64(v)
	case uint64:
		v2 = uint64(v)
	case float32:
		v2 = uint64(v)
	case float64:
		v2 = uint64(v)
	case bool:
		if v {
			v2 = 1
		} else {
			v2 = 0
		}
	case string:
		v3, err := strconv.ParseFloat(v, 64)
		if err == nil {
			v2 = uint64(v3)
		} else {
			v2 = 0
		}
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
	return
}

func TypeToInt64(v1 any) (v2 int64) {
	switch v := v1.(type) {
	case int:
		v2 = int64(v)
	case int8:
		v2 = int64(v)
	case int16:
		v2 = int64(v)
	case int32:
		v2 = int64(v)
	case int64:
		v2 = int64(v)
	case uint:
		v2 = int64(v)
	case uint8:
		v2 = int64(v)
	case uint16:
		v2 = int64(v)
	case uint32:
		v2 = int64(v)
	case uint64:
		v2 = int64(v)
	case float32:
		v2 = int64(v)
	case float64:
		v2 = int64(v)
	case bool:
		if v {
			v2 = 1
		} else {
			v2 = 0
		}
	case string:
		v3, err := strconv.ParseFloat(v, 64)
		if err == nil {
			v2 = int64(v3)
		} else {
			v2 = 0
		}
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
	return
}

func TypeToFloat64(v1 any) (v2 float64) {
	switch v := v1.(type) {
	case int:
		v2 = float64(v)
	case int8:
		v2 = float64(v)
	case int16:
		v2 = float64(v)
	case int32:
		v2 = float64(v)
	case int64:
		v2 = float64(v)
	case uint:
		v2 = float64(v)
	case uint8:
		v2 = float64(v)
	case uint16:
		v2 = float64(v)
	case uint32:
		v2 = float64(v)
	case uint64:
		v2 = float64(v)
	case float32:
		v2 = float64(v)
	case float64:
		v2 = float64(v)
	case bool:
		if v {
			v2 = 1
		} else {
			v2 = 0
		}
	case string:
		v3, err := strconv.ParseFloat(v, 64)
		if err == nil {
			v2 = v3
		} else {
			v2 = 0
		}
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
	return
}

func TypeToString(v1 any) (v2 string) {
	switch v := v1.(type) {
	case string:
		v2 = v
		return
	case []byte:
		v2 = string(v)
		return
	case []rune:
		v2 = string(v)
		return
	default:
		v2 = fmt.Sprint(v1)
	}
	return
}
