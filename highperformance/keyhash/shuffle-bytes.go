package keyhash

// can be used when you want the result to stay in the same alphabet (UUID, file name, etc)
func ShuffleBytes(key []byte, times int) {
	for t := 0; t < times; t++ {
		for i := 0; i < len(key); i++ {
			v := key[i]
			i2 := (t + int(v)) % len(key)
			key[i], key[i2] = key[i2], key[i]
		}
	}
}

// best for binary data of small length
func ShuffleBytesX(key []byte, times int) {
	for t := 0; t < times; t++ {
		for i := 0; i < len(key); i++ {
			v := key[i]
			i2 := (t + int(v)) % len(key)
			key[i], key[i2] = key[i2], key[i]
		}
		for i := 0; i < len(key); i++ {
			i2 := (i + 1) % len(key)
			key[i] ^= key[i2]
		}
	}
}
