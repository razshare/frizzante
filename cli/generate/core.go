package generate

import "path/filepath"

func Core(options CoreOptions) (err error) {
	if err = Copy(CopyOptions{
		From: "internal/project/lib/core",
		To:   filepath.Join("lib", "core"),
		Auto: options.Auto,
		Efs:  options.Efs,
	}); err != nil {
		return err
	}

	if err = Copy(CopyOptions{
		From: "internal/project/app/lib/scripts/core",
		To:   filepath.Join(options.App, "lib", "scripts", "core"),
		Auto: options.Auto,
		Efs:  options.Efs,
	}); err != nil {
		return err
	}

	err = FixImports(FixImportsOptions{Directory: filepath.Join("lib", "core")})

	return
}
