package text

import (
	"crypto/sha1"
	"encoding/base64"
	"strings"
)

func Sha1(txt string) (string, error) {
	hasher := sha1.New()
	_, err := hasher.Write([]byte(txt))
	if err != nil {
		return "", err
	}
	b64 := base64.URLEncoding.EncodeToString(hasher.Sum(nil))
	return strings.ReplaceAll(b64, "=", ""), nil
}
