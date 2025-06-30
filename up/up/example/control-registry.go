package example

/*
	This is a component, to make sure our test examples corrently start up all the components,
	and then correctly wind-up all of them, and that everything is handled exactly once.

	We can't check the order, though. (Because couldn't figure out how). (We could make an int counter,
	and check the order inside the OTDGs, BTW, but the ehole graph isn't clear).
*/

type TControlRegistry struct {
	Components map[string]int
}

var ControlRegistry = &TControlRegistry{
	Components: map[string]int{},
}

func (cr *TControlRegistry) MatchThePattern(patt map[string]int) (ok bool) {
	if len(patt) != len(cr.Components) {
		return false
	}
	for patt_k, patt_v := range patt {
		ok = false
		if v, ok2 := cr.Components[patt_k]; ok2 {
			if v == patt_v {
				ok = true
			}
		}
		if !ok {
			return
		}
	}
	return true
}
