package action

func Config(o ConfigOptions) error {
	err := Generate(GenerateOptions{
		App:      o.App,
		Selected: "bun,air",
		Auto:     o.Auto,
		Platform: o.Platform,
		Go:       o.Go,
		Air:      o.Air,
		Bun:      o.Bun,
		Sqlc:     o.Sqlc,
	})

	if err != nil {
		return err
	}

	return Install(InstallOptions{
		App: o.App,
		Go:  o.Go,
		Bun: o.Bun,
	})
}
