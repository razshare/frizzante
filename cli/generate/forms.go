package generate

import "path/filepath"

func Forms(options FormsOptions) (err error) {
	if err = Copy(CopyOptions{
		From: "internal/additions/app/lib/components/forms",
		To:   filepath.Join("app", "lib", "components", "forms"),
		Auto: options.Auto,
		Efs:  options.Efs,
	}); err != nil {
		return
	}

	return
}
