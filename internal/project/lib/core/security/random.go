package security

import "crypto/rand"

func RandomHex(length int) string {
	bytes := make([]byte, length)
	_, _ = rand.Read(bytes)

	to := make([]byte, length*2)
	var i int
	for _, b := range bytes {
		to[i] = HexTable[b>>4]
		to[i+1] = HexTable[b&0x0f]
		i += 2
	}
	return string(to)
}
