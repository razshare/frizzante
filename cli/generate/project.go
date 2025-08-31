package generate

import (
	"os"
	"path/filepath"
)

func Project(options ProjectOptions) (err error) {
	if err = Copy(CopyOptions{
		Ignore: []string{
			"internal/template/project/lib/database",
			"internal/template/project/lib/session/disk",
			"internal/template/project/app/lib/components/forms",
			"internal/template/project/app/lib/components/links",
		},
		From: "internal/template/project",
		To:   options.Name,
		Auto: options.Auto,
		Efs:  options.Efs,
	}); err != nil {
		return err
	}
	if err = os.Rename(
		filepath.Join(options.Name, "go.mod.txt"),
		filepath.Join(options.Name, "go.mod"),
	); err != nil {
		return
	}

	if err = os.Rename(
		filepath.Join(options.Name, "go.sum.txt"),
		filepath.Join(options.Name, "go.sum"),
	); err != nil {
		return
	}

	return
}
