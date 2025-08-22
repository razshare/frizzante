package codegen

import "path/filepath"

func Project(opts ProjectOptions) error {
	return Copy(CopyOptions{
		From: filepath.Join("template", "project"),
		To:   opts.Name,
		Auto: opts.Auto,
	})
}
