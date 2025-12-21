package generate

import "path/filepath"

func Links(options LinksOptions) (err error) {
	if err = Copy(CopyOptions{
		From: "internal/additions/app/lib/components/links",
		To:   filepath.Join("app", "lib", "components", "links"),
		Efs:  options.Efs,
	}); err != nil {
		return
	}

	return
}
