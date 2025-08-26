package action

import (
	"github.com/razshare/frizzante/cli/codegen"
	"os/exec"
)

func Config(opts ConfigOptions) error {

	if _, err := exec.LookPath(opts.Air); err != nil {
		err = codegen.Air(codegen.AirOptions{
			Air:      opts.Air,
			Auto:     opts.Auto,
			Platform: opts.Platform,
		})

		if err != nil {
			return err
		}
	}

	if _, err := exec.LookPath(opts.Bun); err != nil {
		err = codegen.Bun(codegen.BunOptions{
			Bun:      opts.Bun,
			Auto:     opts.Auto,
			Platform: opts.Platform,
		})

		if err != nil {
			return err
		}
	}

	return Install(InstallOptions{
		App: opts.App,
		Go:  opts.Go,
		Bun: opts.Bun,
	})
}
