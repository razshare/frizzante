package generate

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/razshare/frizzante/internal/project/lib/core/files"
	"github.com/razshare/frizzante/tui/confirm"
	"github.com/razshare/frizzante/tui/messages"
	"github.com/razshare/frizzante/tui/spinner"
)

func Types(options TypesOptions) (err error) {
	var flags []string

	if files.IsFile("types.go") {
		flags = []string{"run", "-tags", "types", "types.go"}
	} else if files.IsFile(filepath.Join("lib", "types.go")) {
		flags = []string{"run", filepath.Join("lib", "types.go")}
	} else if files.IsFile(filepath.Join("lib", "types", "main.go")) {
		flags = []string{"run", filepath.Join("lib", "types", "main.go")}
	} else {
		err = errors.New("could not generate types because neither types.go nor lib/types/main.go were found")
		return
	}

	spin := spinner.New(fmt.Sprintf("running %s %s", options.Go, strings.Join(flags, " ")))

	go spinner.Start(spin)
	if !messages.Command(".", os.Environ(), options.Go, flags...) {
		spinner.Stop(spin)
		err = fmt.Errorf("could not run %s %s", options.Go, strings.Join(flags, " "))
		return
	}
	spinner.Stop(spin)

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

		if !options.Auto {
			if yes, err = confirm.Sendf(true, "%s already exsists. Overwrite?", filepath.Join(options.App, "lib", "types", "gen")); err != nil {
				return
			}

			if !yes {
				messages.Infof("skipping types generation")
				return
			}
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
