package algoutils_test

import (
	"fmt"
	"testing"

	"github.com/rusriver/nutz/algoutils"
	"github.com/stretchr/testify/require"
)

func Test_FastMatchTableTreeValue_01(t *testing.T) {
	var v float64
	var cas int

	table := algoutils.FMTTV[float64, string]{
		"K-01": {true, 5.1, nil},
		"K-02": {true, 17.1, algoutils.FMTTV[float64, string]{
			"aa1": {true, 1.1, nil},
			"aa2": {true, 2.1, nil},
			"aa3": {false, 0, algoutils.FMTTV[float64, string]{
				"sss": {false, 0, algoutils.FMTTV[float64, string]{
					"ddd": {true, 99.5, nil},
				}},
			}},
		}},
	}

	/*
		YAML
		    K-01: [true, 5.1]
		    K-02:
		      - true
		      - 17.1
		      - aa1: [true, 1.1]
		        aa2: [true, 2.1]
		        aa3:
				  - false
		          -
		          - sss:
				    - false
		            -
		            - ddd: [true, 99.5]
	*/

	cas = 1
	fmt.Println()
	v, _ = table.FastMatch([]string{"K-01"})
	fmt.Printf("CASE %v:   %v\n", cas, v)
	require.Equal(t, 5.1, v)
	v = 0

	cas = 2
	fmt.Println()
	v, _ = table.FastMatch([]string{"K-02"})
	fmt.Printf("CASE %v:   %v\n", cas, v)
	require.Equal(t, 17.1, v)
	v = 0

	cas = 3
	fmt.Println()
	v, _ = table.FastMatch([]string{"K-02", "aa3", "sss", "ddd"})
	fmt.Printf("CASE %v:   %v\n", cas, v)
	require.Equal(t, 99.5, v)
	v = 0

	cas = 4
	fmt.Println()
	v, _ = table.FastMatch([]string{"K-02", "aaa", "sss", "ddd", "fff"})
	fmt.Printf("CASE %v:   %v\n", cas, v)
	require.Equal(t, 17.1, v)
	v = 0

	cas = 5
	fmt.Println()
	v, _ = table.FastMatch([]string{"K-02", "aa", "sss"})
	fmt.Printf("CASE %v:   %v\n", cas, v)
	require.Equal(t, 17.1, v)
	v = 0

	cas = 6
	fmt.Println()
	v, _ = table.FastMatch([]string{"K-00", "aa", "sss"})
	fmt.Printf("CASE %v:   %v\n", cas, v)
	require.Equal(t, 0.0, v)
	v = 0

	cas = 7
	fmt.Println()
	v, _ = table.FastMatch([]string{"35cf839bd5b3", "a56c456e", "K-01", "6f7574fb"})
	fmt.Printf("CASE %v:   %v\n", cas, v)
	require.Equal(t, 5.1, v)
	v = 0

	cas = 8
	fmt.Println()
	v, _ = table.FastMatch([]string{"35cf839bd5b3", "a56c456e", "K-02", "6f7574fb"})
	fmt.Printf("CASE %v:   %v\n", cas, v)
	require.Equal(t, 17.1, v)
	v = 0

	cas = 9
	fmt.Println()
	v, _ = table.FastMatch([]string{"K-02", "aa1", "K-02", "aa3"})
	fmt.Printf("CASE %v:   %v\n", cas, v)
	require.Equal(t, 1.1, v)
	v = 0

	cas = 10
	fmt.Println()
	v, _ = table.FastMatch([]string{"K-02", "aa1", "K-02", "aa3", "sss", "ddd"})
	fmt.Printf("CASE %v:   %v\n", cas, v)
	require.Equal(t, 99.5, v)
	v = 0

	cas = 11
	fmt.Println()
	v, _ = table.FastMatch([]string{"K-02", "aa3", "sss", "ddd", "K-02", "aa1"})
	fmt.Printf("CASE %v:   %v\n", cas, v)
	require.Equal(t, 99.5, v)
	v = 0
}
