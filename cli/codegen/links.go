package codegen

import "path/filepath"

func Links(o LinksOptions) error {
	return Copy(CopyOptions{
		From: filepath.Join("template", "project", "app", "frizzante", "link"),
		To:   o.Lib,
		Auto: o.Auto,
	})
}
