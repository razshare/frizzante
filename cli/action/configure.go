package action

import (
	"os/exec"

	"github.com/razshare/frizzante/cli/generate"
	"github.com/razshare/frizzante/internal/project/lib/core/files"
)

func Configure(options ConfigureOptions) (err error) {
	if _, err = exec.LookPath(options.Air); err != nil && !files.IsFile(options.Air) {
		if err = generate.Air(generate.AirOptions{
			Air:      options.Air,
			Auto:     options.Auto,
			Platform: options.Platform,
		}); err != nil {
			return
		}
	}

	if _, err = exec.LookPath(options.Bun); err != nil && !files.IsFile(options.Bun) {
		if err = generate.Bun(generate.BunOptions{
			Bun:      options.Bun,
			Auto:     options.Auto,
			Platform: options.Platform,
		}); err != nil {
			return
		}
	}

	if err = Install(InstallOptions{
		App: options.App,
		Go:  options.Go,
		Bun: options.Bun,
	}); err != nil {
		return
	}

	if err = Package(PackageOptions{
		App:  options.App,
		Bun:  options.Bun,
		Prod: true,
	}); err != nil {
		return
	}

	return
}
