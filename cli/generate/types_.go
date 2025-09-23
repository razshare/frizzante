package generate

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/razshare/frizzante/internal/project/lib/core/files"
	"github.com/razshare/frizzante/tui/confirm"
	"github.com/razshare/frizzante/tui/messages"
	"github.com/razshare/frizzante/tui/spinner"
)

func Types(options TypesOptions) (err error) {
	var types *exec.Cmd
	var spin *spinner.Spinner

	if files.IsFile("types.go") {
		spin = spinner.New("running types.go with -tags types")
		types = exec.Command(options.Go, "run", "-tags", "types", "types.go")
	} else if files.IsFile(filepath.Join("lib", "types", "main.go")) {
		spin = spinner.New("running lib/types/main.go")
		types = exec.Command(options.Go, "run", filepath.Join("lib", "types", "main.go"))
	} else {
		err = errors.New("neither types.go nor lib/types/main.go were found")
		return
	}

	go spinner.Start(spin)
	types.Env = os.Environ()
	// get.Stderr = os.Stderr
	// get.Stdout = os.Stdout
	// get.Stdin = os.Stdin
	err = types.Run()
	spinner.Stop(spin)

	if err != nil {
		if types.Err != nil {
			messages.Error(types.Err.Error())
		}
		return
	}

	if !files.IsDirectory(filepath.Join(".gen", "types")) {
		if err = os.MkdirAll(filepath.Join(".gen", "types"), os.ModePerm); err != nil {
			return
		}
	}

	if !files.IsDirectory(filepath.Join(options.App, "lib", "types")) {
		if err = os.MkdirAll(filepath.Join(options.App, "lib", "types"), os.ModePerm); err != nil {
			return
		}
	}

	// Package "dts" should not be aware of where the "app" directory is,
	// for that reason it generates files directly into ".gen/types".
	// However, the cli knows where the "app" directory is located,
	// so we place those files in the correct directory.

	if files.IsDirectory(filepath.Join(options.App, "lib", "types", "gen")) {
		var yes bool
		if yes, err = confirm.Sendf(true, "%s already exsists. Overwrite?", filepath.Join(options.App, "lib", "types", "gen")); err != nil {
			return
		}

		if !yes {
			messages.Infof("skipping types generation")
			return
		}

		if err = os.RemoveAll(filepath.Join(options.App, "lib", "types", "gen")); err != nil {
			return
		}
	}

	if err = os.Rename(
		filepath.Join(".gen", "types"),
		filepath.Join(options.App, "lib", "types", "gen"),
	); err != nil {
		return
	}

	messages.Successf("types generated in %s", filepath.Join(options.App, "lib", "types", "gen"))

	return
}
