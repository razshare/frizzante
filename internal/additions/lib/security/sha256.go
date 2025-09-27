package security

import "crypto/sha256"

// Sha256 returns the sha256 checksum of the text.
func Sha256(text string) string {
	from := sha256.Sum256([]byte(text))
	to := make([]byte, 64)
	var i int
	for _, b := range from {
		to[i] = HexTable[b>>4]
		to[i+1] = HexTable[b&0x0f]
		i += 2
	}
	return string(to)
}
