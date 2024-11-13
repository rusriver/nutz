package keyhash

// Implementation of Koonin's keyhash 8 bits (one way) and 16 bits (double way) algorithms.

func Get_8bits(key []byte, maxSuffixLength int) (hash int) {

	var h uint8

	var ln = len(key)

	h = uint8(ln)

	for i := 0; i < ln; i++ {
		if maxSuffixLength > 0 && i >= maxSuffixLength {
			break
		}
		h ^= key[ln-i-1]
		h ^= h << 1
	}

	return int(h)
}

func Get_16bits(key []byte, maxSuffixLength int) (hash int) {

	var h, hPrev uint8

	var i1, i2 int = 0, len(key)
	if maxSuffixLength > 0 && i2 > maxSuffixLength {
		i1 = i2 - maxSuffixLength
	}

	h = uint8(i2)

	i2--
	for i1 < i2 {
		hPrev = h
		h ^= key[i1]     // xor
		h ^= key[i2] + 2 // xor, add
		h ^= h << 2      // xor, shift

		i1++
		i2--
	}

	hash = int(hPrev)<<8 + int(h)

	return
}
