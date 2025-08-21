package codegen

import "path/filepath"

func Forms(o FormsOptions) error {
	return Copy(CopyOptions{
		From: filepath.Join("template", "project", "app", "frizzante", "forms"),
		To:   o.Lib,
		Auto: o.Auto,
	})
}
