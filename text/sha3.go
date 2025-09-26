package text

import "golang.org/x/crypto/sha3"

func Sha3_256(text string) string {
	from := sha3.Sum256([]byte(text))
	to := make([]byte, 64)
	var i int
	for _, b := range from {
		to[i] = HexTable[b>>4]
		to[i+1] = HexTable[b&0x0f]
		i += 2
	}
	return string(to)
}

func Sha3_512(text string) string {
	from := sha3.Sum512([]byte(text))
	to := make([]byte, 128)
	var i int
	for _, b := range from {
		to[i] = HexTable[b>>4]
		to[i+1] = HexTable[b&0x0f]
		i += 2
	}
	return string(to)
}
