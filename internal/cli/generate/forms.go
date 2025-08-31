package generate

import "path/filepath"

func Forms(options FormsOptions) error {
	lib := filepath.Join(options.App, "lib")

	return Copy(CopyOptions{
		From: "internal/project/app/lib/components/forms",
		To:   filepath.Join(lib, "components", "forms"),
		Auto: options.Auto,
		Efs:  options.Efs,
	})
}
