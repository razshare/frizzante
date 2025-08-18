package action

import (
	"github.com/razshare/frizzante/files"
	"github.com/razshare/frizzante/tui/messages"
	"os"
	"path/filepath"
)

func Reset(_ ResetOptions) error {
	home := os.Getenv("FRIZZANTE_HOME")
	if home == "" {
		user, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		home = filepath.Join(user, ".frizzante")
	}

	if files.IsDirectory(home) {
		err := os.RemoveAll(home)
		if err != nil {
			return err
		}

		messages.Successf("%s deleted", home)

		return nil
	}

	messages.Infof("%s not found", home)

	return nil
}
