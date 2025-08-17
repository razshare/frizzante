package action

import (
	"github.com/razshare/frizzante/files"
	"github.com/razshare/frizzante/tui/messages"
	"os"
	"path/filepath"
)

func Reset(_ ResetOptions) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	dir := filepath.Join(home, ".frizzante")

	if files.IsDirectory(dir) {
		err = os.RemoveAll(dir)
		if err != nil {
			return err
		}

		messages.Successf("%s deleted", dir)

		return nil
	}

	messages.Infof("%s not found", dir)

	return nil
}
