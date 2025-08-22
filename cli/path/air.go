package path

import (
	"os"
	"path/filepath"
	"strings"
)

func Air(bin string) (string, error) {
	if strings.HasPrefix(bin, "~") {
		dir, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		bin = strings.Replace(bin, "~", dir, 1)
		return bin, nil
	}

	if !strings.Contains(bin, string(filepath.Separator)) {
		return bin, nil
	}

	return bin, nil
}
