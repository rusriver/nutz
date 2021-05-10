package main

import (
	"fmt"
	"strings"
)

type IH struct {
	a  int
	b  string
	c  []int
	d  IJ
	d2 IK
}
type IJ struct {
	z bool
	v bool
}
type IK map[string]bool

func main() {
	ih := map[string]interface{}{
		"a": 0,
		"b": "kjhgkjhg",
		"c": []interface{}{
			1, 2, 3, 4, 5, "fff", 5, 6,
		},
		"d": map[string]interface{}{
			"z": true,
			"x": true,
			"c": false,
			"v": false,
		},
		"d2": map[string]interface{}{
			"a": true,
			"s": true,
			"d": false,
			"f": false,
		},
	}
	fmt.Printf("%+v\n", ih)
	IterlevJ(ih)

	ih2 := IH{}
	fmt.Printf("%+v\n", ih2)
}

var indent = 0

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
