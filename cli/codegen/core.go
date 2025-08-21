package codegen

import "path/filepath"

func Core(o CoreOptions) error {
	return Copy(CopyOptions{
		From: filepath.Join("template", "project", "app", "frizzante", "core"),
		To:   o.Lib,
		Auto: o.Auto,
	})
}
