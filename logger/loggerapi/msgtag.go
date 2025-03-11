package loggerapi

import "strings"

// Few first strings are reported to metrics, be careful to NOT put in them
// high-cardinality IDs.
type Msgtags []string

func (msgtags Msgtags) String() string {
	return strings.Join(msgtags, " ")
}

func (msgtags Msgtags) With(ss ...string) Msgtags {
	if msgtags == nil {
		// msgtag = &Msgtag{ss...} // WTF, doesn't work (TODO)
		return append(Msgtags{}, ss...)
	}
	msgtags = append(msgtags[:0:0], msgtags...)
	for _, s := range ss {
		msgtags = append(msgtags, s)
	}
	return msgtags
}
