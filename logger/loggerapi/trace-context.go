package loggerapi

import "strings"

type TraceContext struct {
	ETC        EventTracingCounters
	Msgtags    Msgtags
	IsInactive bool
}

// faster alternative to the logger.GetTraceContext().Msgtags.CloneWith("*****").String() idiom
func (tc *TraceContext) MsgtagsStringWith(tags ...string) string {
	return tc.Msgtags.String() + " " + strings.Join(tags, " ")
}
