package logger

import "github.com/rusriver/nutz/logger/loggerapi"

func NormalizeTheMsgtag(msgtag *loggerapi.Msgtags, n int) *loggerapi.Msgtags {
	for {
		if len(*msgtag) >= n {
			break
		}
		*msgtag = append(*msgtag, "")
	}
	return msgtag
}
