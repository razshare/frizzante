package generate

import "path/filepath"

func Core(options CoreOptions) (err error) {
	lib := filepath.Join(options.App, "lib")

	if err = Copy(CopyOptions{
		From: "internal/template/project/app/lib/components/core",
		To:   filepath.Join(lib, "components", "core"),
		Auto: options.Auto,
		Efs:  options.Efs,
	}); err != nil {
		return err
	}

	if err = Copy(CopyOptions{
		From: "internal/template/project/app/lib/scripts/core",
		To:   filepath.Join(lib, "scripts", "core"),
		Auto: options.Auto,
		Efs:  options.Efs,
	}); err != nil {
		return err
	}

	return
}
