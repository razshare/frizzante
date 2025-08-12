package path

import (
	"github.com/razshare/frizzante/cli/extension"
	"github.com/razshare/frizzante/cli/state"
	"github.com/razshare/frizzante/tui/messages"
	"os"
	"path/filepath"
	"strings"
)

func Go(base string) string {
	var bin string

	if *state.Go != "" {
		bin = *state.Go
	} else {
		bin = Go(".")
	}

	if strings.HasPrefix(bin, "~") {
		dirname, err := os.UserHomeDir()
		if err != nil {
			messages.Fatal(err)
		}
		bin = strings.Replace(bin, "~", dirname, 1)
		return bin + extension.Find()
	}

	if !strings.Contains(bin, string(filepath.Separator)) {
		return bin + extension.Find()
	}

	var pathError error
	bin, pathError = filepath.Rel(base, bin)
	if pathError != nil {
		messages.Fatal(pathError)
	}

	return bin + extension.Find()
}
