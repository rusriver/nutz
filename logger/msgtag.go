package logger

import "strings"

// Few first strings are reported to metrics, be careful to NOT put in them
// high-cardinality IDs.
type Msgtag []string

func (msgtag *Msgtag) String() string {
	return strings.Join(*msgtag, "-")
}

func (msgtag *Msgtag) With(ss ...string) *Msgtag {
	clone := append((*msgtag)[:0:0], *msgtag...)
	msgtag = &clone
	for _, s := range ss {
		*msgtag = append(*msgtag, s)
	}
	return msgtag
}
