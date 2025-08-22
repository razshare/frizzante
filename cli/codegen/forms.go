package codegen

import "path/filepath"

func Forms(opts FormsOptions) error {
	return Copy(CopyOptions{
		From: filepath.Join("template", "project", "app", "frizzante", "forms"),
		To:   opts.Lib,
		Auto: opts.Auto,
	})
}
