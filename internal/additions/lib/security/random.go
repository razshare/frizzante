package security

import (
	"crypto/rand"
	"encoding/base64"
)

func RandomHex(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	to := make([]byte, length*2)
	var i int
	for _, b := range bytes {
		to[i] = HexTable[b>>4]
		to[i+1] = HexTable[b&0x0f]
		i += 2
	}
	return string(to), nil
}

func RandomBase64(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(bytes), nil
}

func RandomBase64Standard(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(bytes), nil
}
