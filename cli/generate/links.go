package generate

import "path/filepath"

func Links(options LinksOptions) error {
	return Copy(CopyOptions{
		From: filepath.Join("internal", "template", "project", "app", "frizzante", "link"),
		To:   options.Lib,
		Auto: options.Auto,
	})
}
