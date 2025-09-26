package text

import (
	"crypto/rand"
	"encoding/base64"
)

func RandomBytes(length int) ([]byte, error) {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	return bytes, err
}

func RandomHex(length int) (string, error) {
	bytes, err := RandomBytes(length)
	if err != nil {
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
	bytes, err := RandomBytes(length)
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(bytes), nil
}

func RandomBase64Standard(length int) (string, error) {
	bytes, err := RandomBytes(length)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(bytes), nil
}
