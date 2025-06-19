package trace

import (
	"slices"
	"strings"
)

// Can be reported to metrics, be careful to NOT put in them high-cardinality IDs.
// Last elements are most specific, first ones are the least.
type Tags []string

func (tags Tags) String() string {
	return strings.Join(tags, " ")
}

// E.g. to get the one before last, use Last(-2);
// The purpose is to easily get the last tags, which are the most specific;
func (tags Tags) Last(neg_i int) string {
	if neg_i >= 0 {
		return ""
	}
	i := len(tags) + neg_i
	if i >= 0 {
		return tags[i]
	} else {
		return ""
	}
}

func (tags Tags) CloneWith(ss ...string) Tags {
	if tags == nil {
		// msgtag = &Msgtag{ss...} // WTF, doesn't work (TODO)
		return slices.Clone(ss)
	}
	tags = append(tags[:0:0], tags...)
	for _, s := range ss {
		tags = append(tags, s)
	}
	return tags
}
