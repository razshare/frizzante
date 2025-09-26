package text

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/sha512"
)

func HmacSha256(text string, key string) string {
	h := hmac.New(sha256.New, []byte(key))
	h.Write([]byte(text))
	from := h.Sum(nil)
	to := make([]byte, 64)
	var i int
	for _, b := range from {
		to[i] = HexTable[b>>4]
		to[i+1] = HexTable[b&0x0f]
		i += 2
	}
	return string(to)
}

func HmacSha512(text string, key string) string {
	h := hmac.New(sha512.New, []byte(key))
	h.Write([]byte(text))
	from := h.Sum(nil)
	to := make([]byte, 128)
	var i int
	for _, b := range from {
		to[i] = HexTable[b>>4]
		to[i+1] = HexTable[b&0x0f]
		i += 2
	}
	return string(to)
}

func HmacVerifySha256(text string, key string, signature string) bool {
	expectedSignature := HmacSha256(text, key)
	return hmac.Equal([]byte(signature), []byte(expectedSignature))
}

func HmacVerifySha512(text string, key string, signature string) bool {
	expectedSignature := HmacSha512(text, key)
	return hmac.Equal([]byte(signature), []byte(expectedSignature))
}
