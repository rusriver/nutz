package keyhash

func DoubleWrapTheKey(key []byte) []byte {
	var i1, i2 int = 0, len(key)
	i2--
	for i1 < i2 {
		key[i1] ^= key[i2]
		i1++
		i2--
	}
	key = key[:i1]
	return key
}
