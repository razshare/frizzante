package action

func Config(opts ConfigOptions) error {
	err := Generate(GenerateOptions{
		App:      opts.App,
		Selected: "bun,air",
		Auto:     opts.Auto,
		Go:       opts.Go,
		Air:      opts.Air,
		Bun:      opts.Bun,
		Sqlc:     opts.Sqlc,
		Efs:      opts.Efs,
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
