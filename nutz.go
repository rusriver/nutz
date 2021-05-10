package nutz

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

var indent = 0

// err := nutz.MapHierToStruct(hier, &strct)

type tMHTS struct {
	path []string
}

func MapHierToStruct(hier interface{}, strct interface{}) error {
	rv := reflect.ValueOf(strct)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return fmt.Errorf("MapHierToStruct(), must be pointer, not %v", reflect.TypeOf(strct))
	}

	env := &tMHTS{
		path: []string{},
	}
	env.mapHierToStruct(hier, rv)

	return nil
}

func (env *tMHTS) mapHierToStruct(hi1 interface{}, rv1 reflect.Value) error {
	fmt.Printf("MapHierToStruct() at path %v\n", strings.Join(env.path, "."))
EntrySwitch:
	switch hi1v := hi1.(type) {
	default:
		fmt.Printf("+++ EntrySwitch, default\n")
		fmt.Print(strings.Repeat("  ", indent))
		fmt.Printf("%+v : %T\n", hi1, hi1)
	case map[string]interface{}:
		fmt.Printf("+++ EntrySwitch, map[string]interface{}\n")
		var map1struct0 bool

		switch rv1.Kind() {
		case reflect.Interface:
			if rv1.IsNil() {
				rv1.Set(reflect.ValueOf(map[string]interface{}{}))
			} else {
				rv1 = rv1.Elem()
			}
			goto EntrySwitch

		case reflect.Ptr:
			if rv1.IsNil() {
				typ := rv1.Type()
				rv1.Set(reflect.New(typ))
			}
			rv1 = rv1.Elem()
			goto EntrySwitch

		case reflect.Map:
			typ := rv1.Type()
			map_kt := typ.Key().String()
			if map_kt != "string" {
				return fmt.Errorf("MapHierToStruct(), map key type must be [string], not %v, at path %v", map_kt, strings.Join(env.path, "."))
			}
			map1struct0 = true

		case reflect.Struct:
			map1struct0 = false

		default:
			return fmt.Errorf("MapHierToStruct(), incompatible type %v at path %v", rv1.Kind, strings.Join(env.path, "."))
		}

		// if we are here, then there's for sure either a map or struct
		for k, v := range hi1v {
			fmt.Print(strings.Repeat("  ", indent))
			fmt.Printf("%+v => %+v\n", k, v)
			env.path = append(env.path, k)

			// get rv2
			var rv2 reflect.Value
			if map1struct0 { // map
				rv2 = rv1.MapIndex(reflect.ValueOf(k))
				if !rv2.IsValid() || !rv2.CanSet() {
					rv2 = rv1.MapIndex(reflect.ValueOf(ToCamelStr(k)))
					if !rv2.IsValid() || !rv2.CanSet() {
						continue // skip it
					}
				}
			} else { // struct
				rv2 = rv1.FieldByName(k)
				if !rv2.IsValid() || !rv2.CanSet() {
					rv2 = rv1.FieldByName(ToCamelStr(k))
					if !rv2.IsValid() || !rv2.CanSet() {
						continue // skip it
					}
				}
			}

			// go get set it
			var set bool
			switch v2 := v.(type) {
			case map[string]interface{}, []interface{}:
				env.mapHierToStruct(v, rv2)
			case string:
				switch rv2.Kind() {
				case reflect.String:
					rv2.SetString(v2)
					set = true
				case reflect.Slice:
					// the case of byte arrays is not handled intentionally, because of weirdness of Go internals
					if rv2.Type().Elem().String() == "uint8" {
						rv2.SetBytes([]byte(v2))
						set = true
					}
				}
			case bool:
				if rv2.Kind() == reflect.Bool {
					rv2.SetBool(v2)
					set = true
				}
			case uint8, uint16, uint32, uint64:
				var v3 uint64
				switch v2 := v.(type) {
				case uint8:
					v3 = uint64(v2)
				case uint16:
					v3 = uint64(v2)
				case uint32:
					v3 = uint64(v2)
				case uint64:
					v3 = uint64(v2)
				}
				switch rv2.Kind() {
				case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
					rv2.SetUint(v3)
					set = true
				}
			case int8, int16, int32, int64:
				var v3 int64
				switch v2 := v.(type) {
				case int8:
					v3 = int64(v2)
				case int16:
					v3 = int64(v2)
				case int32:
					v3 = int64(v2)
				case int64:
					v3 = int64(v2)
				}
				switch rv2.Kind() {
				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
					rv2.SetInt(v3)
					set = true
				}
			case float32, float64:
				var v3 float64
				switch v2 := v.(type) {
				case float32:
					v3 = float64(v2)
				case float64:
					v3 = float64(v2)
				}
				switch rv2.Kind() {
				case reflect.Float32, reflect.Float64:
					rv2.SetFloat(v3)
					set = true
				}
			default:
				return fmt.Errorf("MapHierToStruct(), unsupported type %T in hierarchy at path %v2", v, strings.Join(env.path, "."))
			}
			if !set {
				return fmt.Errorf("MapHierToStruct(), incompatible type %v at path %v", rv1.Type(), strings.Join(env.path, "."))
			}

			env.path = env.path[:len(env.path)-1]
		}
	case []interface{}:
		fmt.Printf("+++ EntrySwitch, interface{}\n")
		for i, v := range hi1v {
			fmt.Print(strings.Repeat("  ", indent))
			fmt.Printf("[%+v] %+v\n", i, v)
			env.path = append(env.path, strconv.Itoa(i))
			env.mapHierToStruct(v, rv1)
			env.path = env.path[:len(env.path)-1]
		}
	}

	return nil
}

func IterlevJ(rnode interface{}) {
	switch rnv := rnode.(type) {
	default:
		fmt.Print(strings.Repeat("  ", indent))
		fmt.Printf("%+v : %T\n", rnode, rnode)
	case map[string]interface{}:
		for k, v := range rnv {
			fmt.Print(strings.Repeat("  ", indent))
			fmt.Printf("%+v => %+v\n", k, v)
			indent++
			IterlevJ(v)
			indent--
		}
	case []interface{}:
		for i, v := range rnv {
			fmt.Print(strings.Repeat("  ", indent))
			fmt.Printf("[%+v] %+v\n", i, v)
			indent++
			IterlevJ(v)
			indent--
		}
	}
}
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
