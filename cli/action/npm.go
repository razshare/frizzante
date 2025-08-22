package action

import (
	"github.com/razshare/frizzante/cli/npm"
	"github.com/razshare/frizzante/tui/npmselect"
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

	bun, err := filepath.Rel(opts.App, opts.Bun)
	if err != nil {
		return err
	}

	return npm.Install(bun, opts.App, pkgs...)
}
