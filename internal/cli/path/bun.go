package path

import (
	"os"
	"strings"
)

func Bun(bin string) (out string, err error) {
	if strings.HasPrefix(bin, "~") {
		var home string
		if home, err = os.UserHomeDir(); err != nil {
			return
		}
		out = strings.Replace(bin, "~", home, 1)
		return
	}
	out = bin
	return
}
