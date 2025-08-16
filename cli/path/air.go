package path

import (
	"github.com/razshare/frizzante/cli/extension"
	"github.com/razshare/frizzante/cli/state"
	"os"
	"path/filepath"
	"strings"
)

func Air(base string) (string, error) {
	var bin string

	if *state.Air != "" {
		bin = *state.Air
	} else {
		bin = filepath.Join(".gen", "air", "air")
	}

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

	bin, err := filepath.Rel(base, bin)
	if err != nil {
		return "", err
	}

	return bin + extension.Find(), nil
}
