package generate

import "path/filepath"

func TypesFeature(options TypesFeatureOptions) (err error) {
	if err = Copy(CopyOptions{
		From: "internal/additions/lib/types",
		To:   filepath.Join("lib", "types"),
		Auto: options.Auto,
		Efs:  options.Efs,
	}); err != nil {
		return
	}

	err = FixImports(FixImportsOptions{Directory: filepath.Join("lib", "types")})

	return
}
