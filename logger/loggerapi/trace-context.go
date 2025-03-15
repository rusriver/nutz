package loggerapi

type TraceContext struct {
	ETC        EventTracingCounters
	Msgtags    Msgtags
	IsInactive bool
}
