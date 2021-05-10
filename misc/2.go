package main

import (
	"fmt"
	"reflect"
)

type person struct {
	firsName string
	lastName string
	iceCream []string
}

func main() {
	u := struct {
		myMap    map[int]int
		mySlice  []string
		myPerson person
	}{
		myMap:   map[int]int{1: 10, 2: 20},
		mySlice: []string{"red", "green"},
		myPerson: person{
			firsName: "Esmaeil",
			lastName: "Abedi",
			iceCream: []string{"Vanilla", "chocolate"},
		},
	}
	v := reflect.ValueOf(u)
	for i := 0; i < v.NumField(); i++ {
		fmt.Println(v.Type().Field(i).Name)
		fmt.Println("\t", v.Field(i))
	}
}
