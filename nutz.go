package nutz

import (
	"fmt"
	"reflect"
	"strings"
)

var indent = 0

func IterlevS(v reflect.Value) {
	switch v.Kind() {
	case reflect.Interface, reflect.Ptr:
		if v.IsNil() {
			fmt.Print(strings.Repeat("   ", indent+1))
			fmt.Printf("nil\n")
			return
		}
		v2 := v.Elem()
		fmt.Print(strings.Repeat("   ", indent))
		fmt.Printf("dereference, ptr or interface\n")
		if !v2.CanInterface() {
			fmt.Print(strings.Repeat("   ", indent+1))
			fmt.Printf("unexported\n")
		}
		indent++
		IterlevS(v2)
		indent--
	case reflect.Map:
		if v.IsNil() {
			fmt.Print(strings.Repeat("   ", indent+1))
			fmt.Printf("nil\n")
			return
		}
		fmt.Print(strings.Repeat("   ", indent))
		fmt.Printf("key type: %v\n", v.Type().Key())
		for _, key := range v.MapKeys() {
			v2 := v.MapIndex(key)
			fmt.Print(strings.Repeat("   ", indent))
			fmt.Printf("k %v %v : %v, %v\n", key, v2, v2.Kind(), v2.Type())
			if v2.CanInterface() {
				indent++
				IterlevS(v2)
				indent--
			} else {
				fmt.Print(strings.Repeat("   ", indent+1))
				fmt.Printf("unexported\n")
			}
		}
	case reflect.Slice:
		if v.IsNil() {
			fmt.Print(strings.Repeat("   ", indent+1))
			fmt.Printf("nil\n")
			return
		}
		fallthrough
	case reflect.Array:
		for i := 0; i < v.Len(); i++ {
			v2 := v.Index(i)
			fmt.Print(strings.Repeat("   ", indent))
			fmt.Printf("i %v %v : %v, %v\n", i, v2, v2.Kind(), v2.Type())
			indent++
			IterlevS(v2)
			indent--
		}
	case reflect.String:
		str := v.Interface().(string)
		fmt.Print(strings.Repeat("   ", indent))
		fmt.Printf("%v : %T\n", str, str)
	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			v2 := v.Field(i)
			key := v.Type().Field(i).Name
			fmt.Print(strings.Repeat("   ", indent))
			fmt.Printf("sf %v %v : %v, %v\n", key, v2, v2.Kind(), v2.Type())
			if v2.CanInterface() {
				indent++
				IterlevS(v2)
				indent--
			} else {
				fmt.Print(strings.Repeat("   ", indent+1))
				fmt.Printf("unexported\n")
			}
		}
	}
}

func ToCamelStr(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return s
	}

	n := strings.Builder{}
	n.Grow(len(s))
	capNext := true
	for _, v := range []byte(s) {
		vIsCap := v >= 'A' && v <= 'Z'
		vIsLow := v >= 'a' && v <= 'z'
		if capNext {
			if vIsLow {
				v -= 32
			}
		}
		if vIsCap || vIsLow {
			n.WriteByte(v)
			capNext = false
		} else if vIsNum := v >= '0' && v <= '9'; vIsNum {
			n.WriteByte(v)
			capNext = true
		} else {
			capNext = v == '_' || v == ' ' || v == '-' || v == '.'
		}
	}
	return n.String()
}
