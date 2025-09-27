package security

import "golang.org/x/crypto/sha3"

// Sha3Variant256 returns the sha3-256 digest of the text.
func Sha3Variant256(text string) string {
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

// Sha3Variant512 returns the sha3-512 digest of the text.
func Sha3Variant512(text string) string {
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
