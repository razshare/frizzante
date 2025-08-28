package text

import (
	"crypto/sha1"
	"encoding/base64"
	"strings"
)

func Sha1(text string) (hash string, err error) {
	hasher := sha1.New()
	if _, err = hasher.Write([]byte(text)); err != nil {
		return
	}
	text64 := base64.URLEncoding.EncodeToString(hasher.Sum(nil))
	hash = strings.ReplaceAll(text64, "=", "")
	return
}
