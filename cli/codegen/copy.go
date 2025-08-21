package codegen

import (
	"fmt"
	"github.com/razshare/frizzante/cli/user"
	"github.com/razshare/frizzante/files"
	"github.com/razshare/frizzante/tui/confirm"
	"github.com/razshare/frizzante/tui/messages"
	"os"
	"path/filepath"
)

func Copy(o CopyOptions) error {
	if files.IsFile(o.To) || files.IsDirectory(o.To) {
		if !o.Auto {
			yes, err := confirm.Sendf(true, "%s already exists. Overwrite?", o.To)
			if err != nil {
				return err
			}

			if !yes {
				messages.Infof("skipping %s", o.To)
				return nil
			}
		}

		err := os.RemoveAll(o.To)
		if err != nil {
			return err
		}
	}

	home, err := user.FrizzanteHome()
	if err != nil {
		return err
	}

	if files.IsDirectory(filepath.Join(home, o.From)) {
		err = files.CopyDirectory(filepath.Join(home, o.From), o.To)
		if err != nil {
			return err
		}
	} else if files.IsFile(filepath.Join(home, o.From)) {
		err = files.CopyFile(filepath.Join(home, o.From), o.To)
		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("%s not found", o.From)
	}

	messages.Successf("%s created", o.To)

	return nil
}
