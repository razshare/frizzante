package codegen

import "path/filepath"

func Links(opts LinksOptions) error {
	return Copy(CopyOptions{
		From: filepath.Join("internal", "template", "project", "app", "frizzante", "link"),
		To:   opts.Lib,
		Auto: opts.Auto,
	})
}
