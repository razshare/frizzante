package path

import (
	"github.com/razshare/frizzante/cli/extension"
	"os"
	"path/filepath"
	"strings"
)

func Air(bin string) (string, error) {
	if strings.HasPrefix(bin, "~") {
		dirname, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		bin = strings.Replace(bin, "~", dirname, 1)
		return bin + extension.Find(), nil
	}

	if !strings.Contains(bin, string(filepath.Separator)) {
		return bin + extension.Find(), nil
	}

	return bin + extension.Find(), nil
}
