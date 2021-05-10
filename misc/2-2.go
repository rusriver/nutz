package main

import (
	"fmt"
	"reflect"
	"strings"
)

type person struct {
	firsName string
	LastName string
	iceCream []string
}

type someT int

func main() {
	u := struct {
		MyMap     map[int]int
		MyMap2    map[string]interface{}
		MySlice   []string
		MySlice2  []string
		MyPerson  person
		MyPerson2 *person
		MyPerson3 *person
		age       int
		Some      someT
	}{
		MyMap:   map[int]int{1: 10, 2: 20},
		MyMap2:  map[string]interface{}{"a": 10, "A": "AAA-string"},
		MySlice: []string{"red", "green"},
		MyPerson: person{
			firsName: "Esmaeil",
			LastName: "Abedi",
			iceCream: []string{"Vanilla", "chocolate"},
		},
		MyPerson3: &person{
			LastName: "Maya",
		},
		age:  15,
		Some: 25,
	}

	IterlevS(reflect.ValueOf(u))
}

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
