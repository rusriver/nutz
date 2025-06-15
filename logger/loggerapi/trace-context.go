package loggerapi

import "strings"

type TraceContext struct {
	ETC        EventTracingCounters // []int
	Msgtags    Msgtags              // []string
	IsInactive bool
}

// faster alternative to the logger.GetTraceContext().Msgtags.CloneWith("*****").String() idiom
func (tc *TraceContext) MsgtagsStringWith(tags ...string) string {
	return tc.Msgtags.String() + " " + strings.Join(tags, " ")
}
