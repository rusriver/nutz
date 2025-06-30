package trace

import "strings"

type Forwardtrace struct {
	ETC  EventTracingCounters // []int
	Tags Tags                 // []string
}

// faster alternative to the ...CloneWith("*****").String() idiom
func (ft *Forwardtrace) TagsStringWith(tags ...string) string {
	return ft.Tags.String() + " " + strings.Join(tags, " ")
}
