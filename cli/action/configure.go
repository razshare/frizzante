package action

import (
	"github.com/razshare/frizzante/cli/generate"
	"github.com/razshare/frizzante/files"
	"os/exec"
)

func Configure(options ConfigureOptions) (err error) {
	if _, err = exec.LookPath(options.Air); err != nil && !files.IsFile(options.Air) {
		if err = generate.Air(generate.AirOptions{Air: options.Air, Auto: options.Auto, Platform: options.Platform}); err != nil {
			return
		}
	}

	if _, err = exec.LookPath(options.Bun); err != nil && !files.IsFile(options.Bun) {
		if err = generate.Bun(generate.BunOptions{Bun: options.Bun, Auto: options.Auto, Platform: options.Platform}); err != nil {
			return
		}
	}

	return Install(InstallOptions{
		App: options.App,
		Go:  options.Go,
		Bun: options.Bun,
	})
}
