package generate

import (
	"os"
	"os/exec"
	"path/filepath"

	"github.com/razshare/frizzante/internal/project/lib/core/files"
	"github.com/razshare/frizzante/tui/confirm"
	"github.com/razshare/frizzante/tui/messages"
	"github.com/razshare/frizzante/tui/spinner"
)

func TypesDefinitions(options DefinitionsOptions) (err error) {
	spin := spinner.New("running program with -tags dry,types")

	go spinner.Start(spin)
	get := exec.Command(options.Go, "run", "-tags", "dry,types", "main.go")
	get.Env = append(os.Environ())
	get.Stderr = os.Stderr
	get.Stdout = os.Stdout
	get.Stdin = os.Stdin
	err = get.Run()
	spinner.Stop(spin)

	if err != nil {
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
