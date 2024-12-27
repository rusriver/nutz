package lookup

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLookup_01(t *testing.T) {
	_, err := P(1, "a")
	assert.Equal(t, err.Error(), "not found")
}

func TestLookup_Map(t *testing.T) {
	value, err := P(map[string]int{"foo": 42}, "foo")
	assert.Nil(t, err)
	assert.Equal(t, value.Int(), int64(42))
}

func TestLookup_Ptr(t *testing.T) {
	value, err := P(&structFixture, "String")
	assert.Nil(t, err)
	assert.Equal(t, value.String(), "foo")
}

func TestLookup_Interface(t *testing.T) {
	value, err := P(structFixture, "Interface")

	assert.Nil(t, err)
	assert.Equal(t, value.String(), "foo")
}

func TestLookup_StructBasic(t *testing.T) {
	value, err := P(structFixture, "String")
	assert.Nil(t, err)
	assert.Equal(t, value.String(), "foo")
}

func TestLookup_StructPlusMap(t *testing.T) {
	value, err := P(structFixture, "Map", "foo")
	assert.Nil(t, err)
	assert.Equal(t, value.Int(), int64(42))
}

func TestLookup_MapNamed(t *testing.T) {
	value, err := P(mapFixtureNamed, "foo")
	assert.Nil(t, err)
	assert.Equal(t, value.Int(), int64(42))
}

func TestLookup_NotFound(t *testing.T) {
	_, err := P(structFixture, "qux")
	assert.Equal(t, err, ErrKeyNotFound)

	_, err = P(mapFixture, "qux")
	assert.Equal(t, err, ErrKeyNotFound)
}

func TestLookup_StructIndex(t *testing.T) {
	value, err := P(structFixture, "StructSlice", "0", "Map", "foo")

	assert.Nil(t, err)
	assert.EqualValues(t, value.Interface(), 42)
}

func TestLookup_StructNestedMap(t *testing.T) {
	value, err := P(structFixture, "StructSlice", "0", "String")

	assert.Nil(t, err)
	assert.EqualValues(t, value.Interface(), "foo")
}

func TestLookup_StructNested(t *testing.T) {
	value, err := P(structFixture, "StructSlice", "1", "StructSlice", "1", "String")

	assert.Nil(t, err)
	assert.EqualValues(t, "baz", value.Interface())
}

func TestLookupString_Complex(t *testing.T) {
	value, err := DotP(structFixture, "StructSlice.0.Map.foo")
	assert.Nil(t, err)
	assert.EqualValues(t, value.Interface(), 42)

	value, err = DotP(mapComplexFixture, "map.bar")
	assert.Nil(t, err)
	assert.EqualValues(t, value.Interface(), 1)
}

func TestLookup_EmptySlice(t *testing.T) {
	fixture := [][]MyStruct{{}}
	_, err := DotP(fixture, "String")
	assert.Equal(t, err, ErrInvalidIndexUsage)
}

func TestLookup(t *testing.T) {
	_, err := P(structFixture, "String-Missing")
	assert.Equal(t, err, ErrKeyNotFound)
}

func TestLookup_Map2(t *testing.T) {
	_, err := P(map[string]int{"Foo": 42}, "foo")
	assert.Equal(t, err, ErrKeyNotFound)
}

func TestLookup_ListPtr(t *testing.T) {
	type Inner struct {
		Value string
	}

	type Outer struct {
		Values *[]Inner
	}

	values := []Inner{{Value: "first"}, {Value: "second"}}
	data := Outer{Values: &values}

	value, err := DotP(data, "Values.0.Value")
	assert.Nil(t, err)
	assert.Equal(t, value.String(), "first")
}

func TestLookup_Ptr_Index(t *testing.T) {
	ptr := &structFixture
	value, err := DotP(ptr, "StructSlice.1.String")
	assert.Nil(t, err)
	assert.Equal(t, value.String(), "qux")
}

func TestLookup_IndexOutOfBounds(t *testing.T) {
	_, err := DotP(structFixture, "StructSlice.42.String")
	assert.Equal(t, err, ErrIndexOutOfBounds)
}

func ExampleDotP() {
	type Cast struct {
		Actor, Role string
	}

	type Serie struct {
		Cast []Cast
	}

	series := map[string]Serie{
		"A-Team": {Cast: []Cast{
			{Actor: "George Peppard", Role: "Hannibal"},
			{Actor: "Dwight Schultz", Role: "Murdock"},
			{Actor: "Mr. T", Role: "Baracus"},
			{Actor: "Dirk Benedict", Role: "Faceman"},
		}},
	}

	q := "A-Team.Cast.0.Actor"
	value, _ := DotP(series, q)
	fmt.Println(q, "->", value.Interface())

	// Output:
	// A-Team.Cast.0.Actor -> George Peppard
}

func ExampleP() {
	type ExampleStruct struct {
		Values struct {
			Foo int
		}
	}

	i := ExampleStruct{}
	i.Values.Foo = 10

	value, _ := P(i, "Values", "Foo")
	fmt.Println(value.Interface())
	// Output: 10
}
