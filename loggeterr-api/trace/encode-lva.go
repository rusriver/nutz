package trace

// See GetIdsChainLVA(), used in pair
func EncodeLVA(ids []string) []byte {
	buf := make([]byte, 0, 0xFF)
	for _, id := range ids {
		buf = append(buf, uint8(len(id)))
		buf = append(buf, []byte(id)...)
	}
	return buf

}
