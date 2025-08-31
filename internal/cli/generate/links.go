package generate

import "path/filepath"

func Links(options LinksOptions) error {
	lib := filepath.Join(options.App, "lib")

	return Copy(CopyOptions{
		From: "internal/template/project/app/lib/components/links",
		To:   filepath.Join(lib, "components", "links"),
		Auto: options.Auto,
		Efs:  options.Efs,
	})
}
