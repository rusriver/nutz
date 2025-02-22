package factory

// See https://en.wikipedia.org/wiki/Medical_record
type ProductProgressRecord[PRODUCT any] struct {
	Product         PRODUCT
	CurrentPosition struct {
		Line    string
		Station int
	}
	FullStopNow bool
	JumpToLine  string
}

func NewProductProgressRecord[PRODUCT any](product PRODUCT, lineName string) *ProductProgressRecord[PRODUCT] {
	pprec := &ProductProgressRecord[PRODUCT]{
		Product: product,
	}
	pprec.CurrentPosition.Line = lineName
	return pprec
}
