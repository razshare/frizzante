package generate

import "path/filepath"

func Forms(options FormsOptions) (err error) {
	if err = Copy(CopyOptions{
		From: "internal/additions/app/lib/components/forms",
		To:   filepath.Join("app", "lib", "components", "forms"),
		Efs:  options.Efs,
	}); err != nil {
		return
	}

	return
}
