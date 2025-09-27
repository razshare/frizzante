package security

import (
	"crypto/rand"
	"encoding/base64"
)

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

func RandomBase64(length int) string {
	bytes := make([]byte, length)
	_, _ = rand.Read(bytes)
	return base64.RawURLEncoding.EncodeToString(bytes)
}

func RandomBase64Standard(length int) string {
	bytes := make([]byte, length)
	_, _ = rand.Read(bytes)
	return base64.StdEncoding.EncodeToString(bytes)
}
