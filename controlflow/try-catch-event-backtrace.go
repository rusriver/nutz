package controlflow

import "github.com/rusriver/nutz/loggeterr-api/trace"

// First element is most specific, last ones are least specific - this is the inverse of forwardtrace tags;
type FullBacktrace []*BacktraceDomainChain

// First element is most specific, last ones are least specific - this is the inverse of forwardtrace tags;
type BacktraceDomainChain []*trace.BacktracePoint

// Used for efficient matching of Id sequences using bytes.Contains() and
// bytes.HasPrefix(). See also trace.EncodeLVA()
func (b FullBacktrace) GetIdsChainLVA() []byte {
	buf := make([]byte, 0, 0xFF)
	for _, d := range b {
		for _, bp := range *d {
			buf = append(buf, uint8(len(bp.Id)))
			buf = append(buf, []byte(bp.Id)...)
		}
	}
	return buf
}

// Used for efficient matching of Id sequences using bytes.Contains() and
// bytes.HasPrefix(). See also trace.EncodeLVA()
func (d *BacktraceDomainChain) GetIdsChainLVA() []byte {
	buf := make([]byte, 0, 0xFF)
	for _, bp := range *d {
		buf = append(buf, uint8(len(bp.Id)))
		buf = append(buf, []byte(bp.Id)...)
	}
	return buf
}

// E.g. to get the one before last, use Last(-2);
// The purpose is to easily get the last thing. Also used internally
func (b FullBacktrace) Last(neg_i int) *BacktraceDomainChain {
	none := &BacktraceDomainChain{trace.BacktracePointf("4319f4a0f592", "non-existing backtrace chain")}
	if neg_i >= 0 {
		return none
	}
	i := len(b) + neg_i
	if i >= 0 {
		return b[i]
	} else {
		return none
	}
}

// E.g. to get the one before last, use Last(-2);
// The purpose is to easily get the last thing. Also used internally
func (d BacktraceDomainChain) Last(neg_i int) *trace.BacktracePoint {
	if neg_i >= 0 {
		return nil
	}
	i := len(d) + neg_i
	if i >= 0 {
		return d[i]
	} else {
		return nil
	}
}
