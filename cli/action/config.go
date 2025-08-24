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

	return codegen.Bun(codegen.BunOptions{
		Bun:      opts.Bun,
		Auto:     opts.Auto,
		Platform: opts.Platform,
	})
}
