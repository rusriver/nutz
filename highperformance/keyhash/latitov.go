package keyhash

func UniformDistributionHash16(key []byte) int {
	var hash, h uint16
	var s01 bool
	var lenKey = len(key)
	for i := 0; i < lenKey; i++ {
		if !s01 {
			h = uint16(key[i]) << 8
			i2 := (i << 3) * int(h)
			h += uint16(i2+(lenKey)) << 8
		} else {
			h += uint16(key[i])
			i2 := (i << 3) * int(h)
			h += uint16(i2 + (lenKey))
			// getHash_applyNext(&hash, h)
			hash += h + (h << 3)

		}
		s01 = !s01
	}
	if s01 {
		// getHash_applyNext(&hash, h)
		hash += h + (h << 3)
	}
	hash ^= (hash & 0xFF00) >> 8
	return int(hash)
}

func getHash_applyNext(hash *uint16, h uint16) {
	(*hash) += h + (h << 3)
	return
}
