package actions

import (
	"github.com/razshare/frizzante/cli/npm"
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
	return npm.Install(npm.InstallOptions{
		Bun:          options.Bun,
		PackageNames: packages,
	})
}
