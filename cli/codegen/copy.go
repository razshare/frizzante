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

func Copy(options CopyOptions) (err error) {
	if files.IsFile(options.To) || files.IsDirectory(options.To) {
		if !options.Auto {
			var overwrite bool
			if overwrite, err = confirm.Sendf(true, "%s already exists. Overwrite?", options.To); err != nil {
				return
			}

			if !overwrite {
				messages.Infof("skipping %s", options.To)
				return
			}
		}

		if err = os.RemoveAll(options.To); err != nil {
			return
		}
	}

	var home string
	if home, err = user.FrizzanteHome(); err != nil {
		return
	}

	if files.IsDirectory(filepath.Join(home, options.From)) {
		if err = files.CopyDirectory(filepath.Join(home, options.From), options.To); err != nil {
			return
		}
	} else if files.IsFile(filepath.Join(home, options.From)) {
		if err = files.CopyFile(filepath.Join(home, options.From), options.To); err != nil {
			return
		}
	} else {
		err = fmt.Errorf("%s not found", options.From)
		return
	}

	messages.Successf("%s created", options.To)

	return
}
