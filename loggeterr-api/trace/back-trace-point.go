package trace

import "fmt"

// Also used as an origin (cause) of the event
type BacktracePoint struct {
	Id     string
	Detail any
}

// new backtrace point formatted
func BacktracePointf(id string, format string, args ...any) *BacktracePoint {
	return &BacktracePoint{
		Id:     id,
		Detail: fmt.Sprintf(format, args...),
	}
}
