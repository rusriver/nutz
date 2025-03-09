package loggerapi

import "strings"

// Few first strings are reported to metrics, be careful to NOT put in them
// high-cardinality IDs.
type Msgtags []string

func (msgtag *Msgtags) String() string {
	return strings.Join(*msgtag, " ")
}

func (msgtag *Msgtags) With(ss ...string) *Msgtags {
	if msgtag == nil {
		// msgtag = &Msgtag{ss...} // WTF, doesn't work (TODO)
		msgtag_v := append(Msgtags{}, ss...)
		return &msgtag_v
	}
	clone := append((*msgtag)[:0:0], *msgtag...)
	msgtag = &clone
	for _, s := range ss {
		*msgtag = append(*msgtag, s)
	}
	return msgtag
}
