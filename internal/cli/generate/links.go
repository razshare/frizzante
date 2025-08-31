package generate

import "path/filepath"

func Links(options LinksOptions) error {
	return Copy(CopyOptions{
		From: "app/lib/components/links",
		To:   filepath.Join(options.App, "lib", "components", "links"),
		Auto: options.Auto,
		Efs:  options.Efs,
	})
}
