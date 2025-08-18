package path

import (
	"os"
	"path/filepath"
	"strings"
)

func Go(bin string) (string, error) {
	if strings.HasPrefix(bin, "~") {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		bin = strings.Replace(bin, "~", home, 1)
		return bin, nil
	}

	if !strings.Contains(bin, string(filepath.Separator)) {
		return bin, nil
	}

	return bin, nil
}
