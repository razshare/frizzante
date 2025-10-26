package actions

import (
	"os/exec"

	"github.com/razshare/frizzante/cli/generate"
	"github.com/razshare/frizzante/internal/project/lib/core/files"
)

func Configure(options ConfigureOptions) (err error) {
	if _, err = exec.LookPath(options.Air); err != nil || !files.IsFile(options.Air) {
		if aerr := generate.Air(generate.AirOptions{
			Air:      options.Air,
			Auto:     options.Auto,
			Platform: options.Platform,
		}); aerr != nil {
			err = aerr
			return
		}
		err = nil
	}

	if _, err = exec.LookPath(options.Bun); err != nil || !files.IsFile(options.Bun) {
		if berr := generate.Bun(generate.BunOptions{
			Bun:      options.Bun,
			Auto:     options.Auto,
			Platform: options.Platform,
		}); berr != nil {
			err = berr
			return
		}
		err = nil
	}

	return
}
