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

func Air(basepath string) string {
	var air string

	if *flags.Air != "" {
		air = *flags.Air
	} else {
		air = filepath.Join(".gen", "air", "air")
	}

	if strings.HasPrefix(air, "~") {
		dirname, err := os.UserHomeDir()
		if err != nil {
			log.Fatal(err)
		}
		air = strings.Replace(air, "~", dirname, 1)
		return air + extension.Find()
	}

	if !strings.Contains(air, string(filepath.Separator)) {
		return air + extension.Find()
	}

	var pathError error

	air, pathError = filepath.Rel(basepath, air)
	if pathError != nil {
		messages.Fatal(pathError)
	}

	return air + extension.Find()
}
