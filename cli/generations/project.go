package generations

import (
	"os"
	"path/filepath"

	"github.com/razshare/frizzante/internal/project/lib/core/files"
)

func Project(options ProjectOptions) (err error) {
	if err = Copy(CopyOptions{
		From: "internal/project",
		To:   options.Name,
		Efs:  options.Efs,
		Ignore: []string{
			"internal/project/additions",
			"internal/project/app/additions",
		},
	}); err != nil {
		return
	}
	if !files.IsFile(filepath.Join(options.Name, "go.mod")) {
		if err = os.WriteFile(
			filepath.Join(options.Name, "go.mod"),
			[]byte("module main\n\ngo 1.25"),
			os.ModePerm,
		); err != nil {
			return
		}
	}
	err = FixImports(FixImportsOptions{Directory: options.Name})
	return
}
