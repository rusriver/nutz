package controlflow

func Either(ops ...func() bool) (r bool) {
	for _, op := range ops {
		if op() {
			r = true
			break
		}
	}
	return
}

func All(ops ...func() bool) (r bool) {
	r = true
	for _, op := range ops {
		if !op() {
			r = false
			break
		}
	}
	return
}
