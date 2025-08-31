package path

import (
	"os"
	"strings"
)

func Sqlc(bin string) (out string, err error) {
	if strings.HasPrefix(bin, "~") {
		var user string
		if user, err = os.UserHomeDir(); err != nil {
			return "", err
		}
		out = strings.Replace(bin, "~", user, 1)
		return
	}
	out = bin
	return
}
