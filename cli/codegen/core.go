package codegen

import "path/filepath"

func Core(opts CoreOptions) error {
	return Copy(CopyOptions{
		From: filepath.Join("internal", "template", "project", "app", "frizzante", "core"),
		To:   opts.Lib,
		Auto: opts.Auto,
	})
}
