package generate

func Project(options ProjectOptions) (err error) {
	return Copy(CopyOptions{
		Ignore: []string{
			"internal/project/lib/database",
			"internal/project/lib/session/disk",
			"internal/project/app/lib/components/forms",
			"internal/project/app/lib/components/links",
		},
		From: "internal/project",
		To:   options.Name,
		Auto: options.Auto,
		Efs:  options.Efs,
	})
}
