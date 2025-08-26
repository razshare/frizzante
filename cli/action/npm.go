package action

import (
	"github.com/razshare/frizzante/cli/npm"
	"github.com/razshare/frizzante/files"
	"github.com/razshare/frizzante/tui/npmselect"
	"os/exec"
	"path/filepath"
)

func Npm(opts NpmOptions) error {
	pkgs, err := npmselect.Send()
	if err != nil {
		return err
	}

	if len(pkgs) == 0 {
		return nil
	}

	var bun string
	if files.IsFile(opts.Bun) {
		if bun, err = filepath.Rel(opts.App, opts.Bun); err != nil {
			return err
		}
	} else if bun, err = exec.LookPath(opts.Bun); err != nil {
		bun = opts.Bun
	}

	return npm.Install(bun, opts.App, pkgs...)
}
