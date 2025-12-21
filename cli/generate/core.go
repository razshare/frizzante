package generate

import "path/filepath"

func Core(options CoreOptions) (err error) {
	if err = Copy(CopyOptions{
		From: "internal/project/lib/core",
		To:   filepath.Join("lib", "core"),
		Efs:  options.Efs,
	}); err != nil {
		return
	}

	if err = Copy(CopyOptions{
		From: "internal/project/app/lib/scripts/core",
		To:   filepath.Join("app", "lib", "scripts", "core"),
		Efs:  options.Efs,
	}); err != nil {
		return
	}

	if err = Copy(CopyOptions{
		From: "internal/project/app/lib/components/core",
		To:   filepath.Join("app", "lib", "components", "core"),
		Efs:  options.Efs,
	}); err != nil {
		return
	}

	if err = FixImports(FixImportsOptions{Directory: filepath.Join("lib", "core")}); err != nil {
		return
	}

	return
}
