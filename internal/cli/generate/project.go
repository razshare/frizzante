package generate

func Project(options ProjectOptions) (err error) {
	return Copy(CopyOptions{
		Ignore: []string{
			"lib/database",
			"lib/session/disk",
			"app/lib/components/forms",
			"app/lib/components/links",
		},
		From: "",
		To:   options.Name,
		Auto: options.Auto,
		Efs:  options.Efs,
	})
}
