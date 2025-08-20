package action

import (
	"path/filepath"

	"github.com/razshare/frizzante/cli/npm"
	"github.com/razshare/frizzante/tui/npmselect"
)

func Npm(o NpmOptions) error {
	pkgs, err := npmselect.Send()
	if err != nil {
		return err
	}

	if len(pkgs) == 0 {
		return nil
	}

	bun, err := filepath.Rel(o.App, o.Bun)
	if err != nil {
		return err
	}

	return npm.Install(bun, o.App, pkgs...)
}
