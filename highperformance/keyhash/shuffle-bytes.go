package keyhash

func ShuffleBytes(key []byte, times int) []byte {

	for t := 0; t < times; t++ {
		for i := 0; i < len(key); i++ {
			v := key[i]
			i2 := (i + int(v)) % len(key)
			key[i], key[i2] = key[i2], key[i]
		}
	}
	return key
}
