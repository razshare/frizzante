package generate

import "path/filepath"

func Core(options CoreOptions) error {
	return Copy(CopyOptions{
		From: filepath.Join("internal", "template", "project", "app", "frizzante", "core"),
		To:   options.Lib,
		Auto: options.Auto,
	})
}
