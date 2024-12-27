package structvalidate_test

type A struct {
	A  *A
	B  map[string]*A
	V1 string
	V2 int
	V3 any
}

func GetData_01() any {
	data := &A{
		B:  make(map[string]*A),
		V1: "root",
	}
	data.A = data
	data.B["k1"] = &A{V1: "k1-V1"}
	data.B["k2"] = &A{V1: "k2-V1", V2: 15, V3: data}
	data.B["k3"] = &A{V3: []byte("hello 123")}
	data.B["k4"] = &A{V3: []string{"aaa", "sss"}}

	return data
}
