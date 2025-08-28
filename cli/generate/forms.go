package generate

import "path/filepath"

func Forms(options FormsOptions) error {
	return Copy(CopyOptions{
		From: filepath.Join("internal", "template", "project", "app", "frizzante", "forms"),
		To:   options.Lib,
		Auto: options.Auto,
	})
}
