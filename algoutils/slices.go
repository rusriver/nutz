package algoutils

// Augments the https://pkg.go.dev/slices, https://pkg.go.dev/slices#Contains

type Fiber interface {
	string | []byte
}

func HasPrefixSequence_Fibers[T Fiber](s1, s2 []T) bool {
	if len(s2) > len(s1) {
		return false
	}
	for i, v2 := range s2 {
		v1 := s1[i]
		if len(v1) != len(v2) {
			return false
		}
		for j := 0; j < len(v1); j++ {
			if v1[j] != v2[j] {
				return false
			}
		}
	}
	return true
}

func HasPrefixSequence_Comparables[T comparable](s1, s2 []T) bool {
	if len(s2) > len(s1) {
		return false
	}
	for i, v2 := range s2 {
		v1 := s1[i]
		if v1 != v2 {
			return false
		}
	}
	return true
}

func HasSubSequence_Fibers[T Fiber](s1, s2 []T) bool {
	for {
		if len(s1) == 0 {
			return false
		}
		if len(s2) > len(s1) {
			return false
		}
		r := HasPrefixSequence_Fibers(s1, s2)
		if r {
			return true
		}
		s1 = s1[1:]
	}
}

func HasSubSequence_Comparables[T comparable](s1, s2 []T) bool {
	for {
		if len(s1) == 0 {
			return false
		}
		if len(s2) > len(s1) {
			return false
		}
		r := HasPrefixSequence_Comparables(s1, s2)
		if r {
			return true
		}
		s1 = s1[1:]
	}
}
