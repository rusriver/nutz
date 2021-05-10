package main

import (
	"fmt"
	"strings"
)

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
	}

	fmt.Printf("%+v\n", ih)

	iterlev(ih)
}

var indent = 0

func iterlev(rnode interface{}) {
	switch rnv := rnode.(type) {
	default:
		fmt.Print(strings.Repeat("  ", indent))
		fmt.Printf("%+v : %T\n", rnode, rnode)
	case map[string]interface{}:
		for k, v := range rnv {
			fmt.Print(strings.Repeat("  ", indent))
			fmt.Printf("%+v => %+v\n", k, v)
			indent++
			iterlev(v)
			indent--
		}
	case []interface{}:
		for i, v := range rnv {
			fmt.Print(strings.Repeat("  ", indent))
			fmt.Printf("[%+v] %+v\n", i, v)
			indent++
			iterlev(v)
			indent--
		}
	}
}
