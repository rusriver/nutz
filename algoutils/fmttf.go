package algoutils

/*
	The FMTTV and FMTTF are basically copies of each other. Due to the limitations of generics
	in Go, I found no way to implement it, without resorting to using "any", one way or another.
	And with any, the whole idea of generics melts, the static code checks fail to work in IDE.
*/

// Fast-Match Table Tree
type FMTTF[Tv any, Tk1 comparable] map[Tk1]FMTTF_Node[Tv, Tk1]

// Fast-Match Table Tree Node
type FMTTF_Node[Tv any, Tk1 comparable] struct {
	F           func(kk []Tk1) Tv
	FmttfNested FMTTF[Tv, Tk1]
}

type fmttf_matchObject[Tv any, Tk1 comparable] struct {
	v                     func(kk []Tk1) Tv
	sq                    []Tk1
	sq_matched_len_with_V int
}

// Matches the provided sample sequence against the table tree.
// It is fast compared to multiple subsequences search.
// Look at examples in tests, and also at the c59145cfdde048f5a3ef.
func (table FMTTF[Tv, Tk1]) FastMatch(sample_seq []Tk1) (v Tv, ok bool) {
	var best_match_object *fmttf_matchObject[Tv, Tk1]

	// try at all possible entry points of the sample sequence
	for i := 0; i < len(sample_seq); i++ {

		moPtr := &fmttf_matchObject[Tv, Tk1]{}

		// fmt.Println(sample_seq[i])
		table.recurse(sample_seq[i:], moPtr)
		// fmt.Println(*moPtr)

		if moPtr.sq_matched_len_with_V > 0 {
			if best_match_object == nil {
				best_match_object = moPtr
			} else {
				if len(moPtr.sq) > len(best_match_object.sq) {
					best_match_object = moPtr
				}
			}
		}
	}
	if best_match_object != nil {
		if len(best_match_object.sq) > 0 {
			ok = true
			v = best_match_object.v(best_match_object.sq) // we ensure it's not nil at 68738c5a4c0260c180ad
		}
	}
	return
}

// traverse the tree recursively
func (table FMTTF[Tv, Tk1]) recurse(sample_seq []Tk1, moPtr *fmttf_matchObject[Tv, Tk1]) {
	k1 := sample_seq[0]
	restk := sample_seq[1:]

	node, ok := table[k1]
	if ok {
		moPtr.sq = append(moPtr.sq, k1)
		if node.F != nil { // 68738c5a4c0260c180ad
			moPtr.v = node.F
			moPtr.sq_matched_len_with_V = len(moPtr.sq)
		}
		if len(restk) > 0 {
			node.FmttfNested.recurse(restk, moPtr)
		}
	}
}

/*
	Why the version for IFiber is missing? Because, you can't use []byte as a map key.
	And strings are comparable.
*/
