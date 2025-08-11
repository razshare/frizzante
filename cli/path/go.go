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

func Go(basepath string) string {
	var goBinary string

	if *flags.Go != "" {
		goBinary = *flags.Go
	} else {
		goBinary = Go(".")
	}

	if strings.HasPrefix(goBinary, "~") {
		dirname, err := os.UserHomeDir()
		if err != nil {
			log.Fatal(err)
		}
		goBinary = strings.Replace(goBinary, "~", dirname, 1)
		return goBinary + extension.Find()
	}

	if !strings.Contains(goBinary, string(filepath.Separator)) {
		return goBinary + extension.Find()
	}

	var pathError error
	goBinary, pathError = filepath.Rel(basepath, goBinary)
	if pathError != nil {
		messages.Fatal(pathError)
	}

	return goBinary + extension.Find()
}
