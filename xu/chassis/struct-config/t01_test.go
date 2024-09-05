package main

import (
	"encoding/json"
	"fmt"
	"testing"
)

type Config struct {
	Name  string
	Order string
	Def1  int
	Hair  *Hair
	Arr1  []int
}
type Hair struct {
	Color string
}

func Test_01_1(t *testing.T) {
	var cases = [][]byte{
		[]byte(`{"Name": "Platypus", "Order": "Monotremata"}`),
		[]byte(`{"Name": "Platypus", "Order": "Monotremata", "Hair": {}}`),
		[]byte(`{"Name": "Platypus", "Order": "Monotremata", "Hair": {"Color":"red"}}`),
		[]byte(`{"Name": "Platypus", "Order": "Monotremata", "Arr1": [1]}`),
	}

	for _, cas := range cases {

		conf1 := &Config{
			Def1: 333,
			Arr1: []int{1, 2, 3},
		}
		conf1.Hair = &Hair{Color: "black"}

		err := json.Unmarshal(cas, conf1)
		if err != nil {
			fmt.Println("error:", err)
		}

		bb, _ := json.MarshalIndent(conf1, "", "    ")
		fmt.Printf("%s\n", bb)
	}
}
