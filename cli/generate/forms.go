package generate

import "path/filepath"

func Forms(options FormsOptions) (err error) {
	return Copy(CopyOptions{
		From: "internal/additions/app/lib/components/forms",
		To:   filepath.Join(options.App, "lib", "components", "forms"),
		Auto: options.Auto,
		Efs:  options.Efs,
	})
}
