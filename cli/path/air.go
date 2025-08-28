package path

import (
	"os"
	"strings"
)

func Air(bin string) (string, error) {
	if strings.HasPrefix(bin, "~") {
		usr, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		bin = strings.Replace(bin, "~", usr, 1)
		return bin, nil
	}

	return bin, nil
}
