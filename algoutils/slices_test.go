package algoutils_test

import (
	"testing"

	"github.com/rusriver/nutz/algoutils"
	"github.com/stretchr/testify/require"
)

func Test_HasPrefixSequence_01(t *testing.T) {
	s1 := []string{"aaa", "sss", "ddd"}
	r1 := algoutils.HasPrefixSequence_Fibers(s1, []string{"aaa"})
	require.Equal(t, true, r1)

	s2 := [][]byte{{1, 1, 1}, {2, 2, 2}, {3, 3, 3}}
	r2 := algoutils.HasPrefixSequence_Fibers(s2, [][]byte{{1, 1, 1}, {2, 2}})
	require.Equal(t, false, r2)
	r2 = algoutils.HasPrefixSequence_Fibers(s2, [][]byte{{1, 1, 1}, {2, 2, 2}})
	require.Equal(t, true, r2)

	s3 := []int{1, 2, 3, 4}
	r3 := algoutils.HasPrefixSequence_Comparables(s3, []int{1, 2})
	require.Equal(t, true, r3)
	r3 = algoutils.HasPrefixSequence_Comparables(s3, []int{3, 4})
	require.Equal(t, false, r3)

	s4 := []float64{1, 2, 3, 4}
	r4 := algoutils.HasPrefixSequence_Comparables(s4, []float64{1, 2})
	require.Equal(t, true, r4)
}

func Test_HasSubSequence_01(t *testing.T) {
	s1 := []string{"aaa", "sss", "ddd", "fff"}
	r1 := algoutils.HasSubSequence_Fibers(s1, []string{"sss", "ddd"})
	require.Equal(t, true, r1)

	s1_1 := []string{"aaa", "sss", "ddd", "fff"}
	r1_1 := algoutils.HasSubSequence_Fibers(s1_1, []string{"xxx"})
	require.Equal(t, false, r1_1)

	s2 := [][]byte{{1, 1, 1}, {2, 2, 2}, {3, 3, 3}}
	r2 := algoutils.HasSubSequence_Fibers(s2, [][]byte{{2, 2, 2}, {3, 3}})
	require.Equal(t, false, r2)
	r2 = algoutils.HasSubSequence_Fibers(s2, [][]byte{{2, 2, 2}, {3, 3, 3}})
	require.Equal(t, true, r2)

	s3 := []int{1, 2, 3, 4, 5}
	r3 := algoutils.HasSubSequence_Comparables(s3, []int{3, 4})
	require.Equal(t, true, r3)
	r3 = algoutils.HasSubSequence_Comparables(s3, []int{4, 3})
	require.Equal(t, false, r3)

	s4 := []float64{1, 2, 3, 4, 5}
	r4 := algoutils.HasSubSequence_Comparables(s4, []float64{3, 4})
	require.Equal(t, true, r4)
}

func Test_HasSubSequence_02(t *testing.T) {
	s1_1 := []string{"aaa", "sss", "ddd", "fff"}
	r1_1 := algoutils.HasSubSequence_Fibers(s1_1, []string{"xxx"})
	require.Equal(t, false, r1_1)
}

func Test_HasSubSequence_03(t *testing.T) {
	s1_1 := []string{"d3740cd9"}
	r1_1 := algoutils.HasSubSequence_Fibers(s1_1, []string{"d0ed6944"})
	require.Equal(t, false, r1_1)
}
