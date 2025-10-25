package actions

import (
	"os/exec"
	"path/filepath"

	"github.com/razshare/frizzante/cli/npm"
	"github.com/razshare/frizzante/internal/project/lib/core/files"
	"github.com/razshare/frizzante/tui/select_npm_packages"
)

func Npm(options NpmOptions) (err error) {
	var packages []string
	if packages, err = select_npm_packages.Send(); err != nil {
		return
	}

	if len(packages) == 0 {
		return nil
	}

	var bun string
	if files.IsFile(options.Bun) {
		if bun, err = filepath.Rel("app", options.Bun); err != nil {
			return
		}
	} else if bun, err = exec.LookPath(options.Bun); err != nil {
		bun = options.Bun
	}

	return npm.Install(npm.InstallOptions{
		Bun:      bun,
		Packages: packages,
	})
}
