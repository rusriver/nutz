package controlflow

type MatchTable []MatchCase

type MatchCase struct {
	Match func() int // 0 no match, 1 match, -1 break
	Then  func() int
}

func (t MatchTable) Match() {
	for _, cas := range t {
		switch cas.Match() {
		case 1:
			switch cas.Then() {
			case -1:
				return
			}
		case -1:
			return
		}
	}
}
