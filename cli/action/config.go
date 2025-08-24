package action

import "github.com/razshare/frizzante/cli/codegen"

func Config(opts ConfigOptions) error {
	err := codegen.Air(codegen.AirOptions{
		Air:      opts.Air,
		Auto:     opts.Auto,
		Platform: opts.Platform,
	})

	if err != nil {
		return err
	}

	err = codegen.Bun(codegen.BunOptions{
		Bun:      opts.Bun,
		Auto:     opts.Auto,
		Platform: opts.Platform,
	})

	if err != nil {
		return err
	}

	return Install(InstallOptions{
		App: opts.App,
		Go:  opts.Go,
		Bun: opts.Bun,
	})
}
