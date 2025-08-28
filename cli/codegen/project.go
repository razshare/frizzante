package codegen

import "path/filepath"

func Project(options ProjectOptions) error {
	return Copy(CopyOptions{
		From: filepath.Join("internal", "template", "project"),
		To:   options.Name,
		Auto: options.Auto,
	})
}
