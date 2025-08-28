package path

import (
	"os"
	"strings"
)

func Air(bin string) (out string, err error) {
	if strings.HasPrefix(bin, "~") {
		var user string
		if user, err = os.UserHomeDir(); err != nil {
			return
		}
		out = strings.Replace(bin, "~", user, 1)
		return
	}
	out = bin
	return
}
