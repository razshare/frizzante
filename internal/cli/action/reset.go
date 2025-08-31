package action

import (
	"os"
	"path/filepath"

	"github.com/razshare/frizzante/files"
	messages2 "github.com/razshare/frizzante/internal/tui/messages"
)

func Reset(_ ResetOptions) (err error) {
	home := os.Getenv("FRIZZANTE_HOME")
	if home == "" {
		var user string
		if user, err = os.UserHomeDir(); err != nil {
			return
		}
		home = filepath.Join(user, ".frizzante")
	}

	if files.IsDirectory(home) {
		if err = os.RemoveAll(home); err != nil {
			return
		}

		messages2.Successf("%s deleted", home)

		return
	}

	messages2.Infof("%s not found", home)

	return
}
