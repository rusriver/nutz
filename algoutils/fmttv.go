package algoutils

/*
	The FMTTV and FMTTF are basically copies of each other. Due to the limitations of generics
	in Go, I found no way to implement it, without resorting to using "any", one way or another.
	And with any, the whole idea of generics melts, the static code checks fail to work in IDE.
*/

// Fast-Match Table Tree
type FMTTV[Tv any, Tk1 comparable] map[Tk1]FMTTV_Node[Tv, Tk1]

// Fast-Match Table Tree Node
type FMTTV_Node[Tv any, Tk1 comparable] struct {
	UseV        bool
	V           Tv
	FmttvNested FMTTV[Tv, Tk1]
}

type fmttv_matchObject[Tv any, Tk1 comparable] struct {
	v                     Tv
	sq                    []Tk1
	sq_matched_len_with_V int
}

// Matches the provided sample sequence against the table tree.
// It is fast compared to multiple subsequences search.
func (table FMTTV[Tv, Tk1]) FastMatch(sample_seq []Tk1) (v Tv, ok bool) {
	var best_match_object *fmttv_matchObject[Tv, Tk1]

	// try at all possible entry points of the sample sequence
	for i := 0; i < len(sample_seq); i++ {

		moPtr := &fmttv_matchObject[Tv, Tk1]{}

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
			v = best_match_object.v
		}
	}
	return
}

// traverse the tree recursively
func (table FMTTV[Tv, Tk1]) recurse(sample_seq []Tk1, moPtr *fmttv_matchObject[Tv, Tk1]) {
	k1 := sample_seq[0]
	restk := sample_seq[1:]

	node, ok := table[k1]
	if ok {
		moPtr.sq = append(moPtr.sq, k1)
		if node.UseV {
			moPtr.v = node.V
			moPtr.sq_matched_len_with_V = len(moPtr.sq)
		}
		if len(restk) > 0 {
			node.FmttvNested.recurse(restk, moPtr)
		}
	}
}

/*
	Why the version for IFiber is missing? Because, you can't use []byte as a map key.
	And strings are comparable.
*/
