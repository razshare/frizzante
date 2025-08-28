package action

import (
	"github.com/razshare/frizzante/files"
	"github.com/razshare/frizzante/tui/messages"
	"os"
	"path/filepath"
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

		messages.Successf("%s deleted", home)

		return
	}

	messages.Infof("%s not found", home)

	return
}
