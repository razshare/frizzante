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

func Copy(opts CopyOptions) (err error) {
	if files.IsFile(opts.To) || files.IsDirectory(opts.To) {
		if !opts.Auto {
			var overwrite bool
			if overwrite, err = confirm.Sendf(true, "%s already exists. Overwrite?", opts.To); err != nil {
				return
			}

			if !overwrite {
				messages.Infof("skipping %s", opts.To)
				return
			}
		}

		if err = os.RemoveAll(opts.To); err != nil {
			return
		}
	}

	var home string
	if home, err = user.FrizzanteHome(); err != nil {
		return
	}

	if files.IsDirectory(filepath.Join(home, opts.From)) {
		if err = files.CopyDirectory(filepath.Join(home, opts.From), opts.To); err != nil {
			return
		}
	} else if files.IsFile(filepath.Join(home, opts.From)) {
		if err = files.CopyFile(filepath.Join(home, opts.From), opts.To); err != nil {
			return
		}
	} else {
		err = fmt.Errorf("%s not found", opts.From)
		return
	}

	messages.Successf("%s created", opts.To)

	return
}
