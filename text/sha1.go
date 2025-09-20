package text

import (
	"crypto/sha1"
	"encoding/hex"
)

func Sha1(text string) string {
	from := sha1.Sum([]byte(text))
	to := make([]byte, 20)
	for i := len(to) - 1; i >= 0; i-- {
		to[i] = from[i]
	}
	return hex.EncodeToString(to)
}
