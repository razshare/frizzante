package path

import (
	"github.com/razshare/frizzante/cli/extension"
	"github.com/razshare/frizzante/cli/state"
	"os"
	"path/filepath"
	"strings"
)

func Go(base string) (string, error) {
	var err error
	var bin string
	var home string

	if *state.Go != "" {
		bin = *state.Go
	} else {
		bin, err = Go(".")
		if err != nil {
			return "", err
		}
	}

	if strings.HasPrefix(bin, "~") {
		home, err = os.UserHomeDir()
		if err != nil {
			return "", err
		}
		bin = strings.Replace(bin, "~", home, 1)
		return bin + extension.Find(), nil
	}

	if !strings.Contains(bin, string(filepath.Separator)) {
		return bin + extension.Find(), nil
	}

	bin, err = filepath.Rel(base, bin)
	if err != nil {
		return "", err
	}

	return bin + extension.Find(), nil
}
