package generate

import "path/filepath"

func Core(options CoreOptions) (err error) {
	if err = Copy(CopyOptions{
		From: "lib/core",
		To:   filepath.Join("lib", "core"),
		Auto: options.Auto,
		Efs:  options.Efs,
	}); err != nil {
		return err
	}

	if err = Copy(CopyOptions{
		From: "app/lib/scripts/core",
		To:   filepath.Join(options.App, "lib", "scripts", "core"),
		Auto: options.Auto,
		Efs:  options.Efs,
	}); err != nil {
		return err
	}

	return
}
