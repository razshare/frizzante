package path

import (
	"github.com/razshare/frizzante/cli/extension"
	"github.com/razshare/frizzante/cli/flags"
	"github.com/razshare/frizzante/tui/messages"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func Bun(basepath string) string {
	var bun string

	if *flags.Bun != "" {
		bun = *flags.Bun
	} else {
		bun = filepath.Join(".gen", "bun", "bun")
	}

	if strings.HasPrefix(bun, "~") {
		dirname, err := os.UserHomeDir()
		if err != nil {
			log.Fatal(err)
		}
		bun = strings.Replace(bun, "~", dirname, 1)
		return bun + extension.Find()
	}

	if !strings.Contains(bun, string(filepath.Separator)) {
		return bun + extension.Find()
	}

	var pathError error
	bun, pathError = filepath.Rel(basepath, bun)
	if pathError != nil {
		messages.Fatal(pathError)
	}

	return bun + extension.Find()
}
