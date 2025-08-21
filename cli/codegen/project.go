package codegen

import "path/filepath"

func Project(o ProjectOptions) error {
	return Copy(CopyOptions{
		From: filepath.Join("template", "project"),
		To:   o.Name,
		Auto: o.Auto,
	})
}
