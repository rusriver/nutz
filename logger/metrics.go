package logger

func NormalizeTheMsgtag(msgtag *Msgtag, n int) *Msgtag {
	for {
		if len(*msgtag) >= n {
			break
		}
		*msgtag = append(*msgtag, "")
	}
	return msgtag
}
