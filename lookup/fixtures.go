package lookup

type MyStruct struct {
	String      string
	Map         map[string]int
	Nested      *MyStruct
	StructSlice []*MyStruct
	Interface   interface{}
}

type MyKey string

var mapFixtureNamed = map[MyKey]int{"foo": 42}

var mapFixture = map[string]int{"foo": 42}

var caseFixtureStruct = struct {
	Foo       int
	TestField int
	Testfield int
	testField int
}{
	0, 1, 2, 3,
}

var caseFixtureMap = map[string]int{
	"Foo":     0,
	"TestKey": 1,
	"Testkey": 2,
	"testKey": 3,
}

var structFixture = MyStruct{
	String:    "foo",
	Map:       mapFixture,
	Interface: "foo",
	StructSlice: []*MyStruct{
		{Map: mapFixture, String: "foo", StructSlice: []*MyStruct{{String: "bar"}, {String: "foo"}}},
		{Map: mapFixture, String: "qux", StructSlice: []*MyStruct{{String: "qux"}, {String: "baz"}}},
	},
}

var mapComplexFixture = map[string]interface{}{
	"map": map[string]interface{}{
		"bar": 1,
	},
	"list": []map[string]interface{}{
		{"baz": 1},
		{"baz": 2},
		{"baz": 3},
	},
	"sub1": caseFixtureMap,
	"sub2": caseFixtureStruct,
}
