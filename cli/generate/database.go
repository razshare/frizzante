package generate

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/razshare/frizzante/internal/project/lib/core/files"
	"github.com/razshare/frizzante/tui/confirm"
	"github.com/razshare/frizzante/tui/messages"
	"github.com/razshare/frizzante/tui/search"
	"github.com/razshare/frizzante/tui/singleselect"
	"github.com/razshare/frizzante/tui/spinner"
)

func Database(options DatabaseOptions) (err error) {
	if options.Type == "" {
		if options.Auto {
			options.Type = "sqlite"
		} else if options.Type, err = singleselect.Send([]search.Choice{{Id: "sqlite"}}, "what type of database would you like to setup?"); err != nil {
			return
		}
	}

	if options.Type != "sqlite" {
		err = fmt.Errorf("database of type %s id not supported", options.Type)
		return
	}

	dbtype := strings.ToLower(options.Type)

	lib := filepath.Join("lib", "database", dbtype)

	if files.IsDirectory(lib) {
		if !options.Auto {
			var overwrite bool
			if overwrite, err = confirm.Sendf(true, "%s already exists. Overwrite?", lib); err != nil {
				return
			}

			if !overwrite {
				messages.Infof("skipping %s", lib)
				return nil
			}
		}

		if err = os.RemoveAll(lib); err != nil {
			return
		}
	}

	if err = Copy(CopyOptions{
		From: "internal/additions/lib/database/" + dbtype,
		To:   lib,
		Auto: options.Auto,
		Efs:  options.Efs,
	}); err != nil {
		return
	}

	if dbtype == "sqlite" {
		spin := spinner.New("adding github.com/mattn/go-sqlite3")

		go spinner.Start(spin)
		install := exec.Command(options.Go, "get", "github.com/mattn/go-sqlite3")
		install.Env = append(os.Environ())
		install.Stderr = os.Stderr
		install.Stdout = os.Stdout
		install.Stdin = os.Stdin
		err = install.Run()
		spinner.Stop(spin)

		if err != nil {
			return
		}

		spin = spinner.New("updating go dependencies")

		go spinner.Start(spin)
		get := exec.Command(options.Go, "get", "-u", "./...")
		get.Env = append(os.Environ())
		get.Stderr = os.Stderr
		get.Stdout = os.Stdout
		get.Stdin = os.Stdin
		err = get.Run()
		spinner.Stop(spin)

		if err != nil {
			return
		}

		messages.Success("sqlite database is ready")

		if strings.Contains(strings.ToLower(options.Generate), "queries") {
			var queries bool
			if queries, err = confirm.Send(true, "would you like to also generate your queries?"); err != nil {
				return
			}

			yaml := "lib/core/database/sqlite/sqlc.yaml"

			//if choice == "mysql" {
			//	yaml = "lib/core/database/mysql/sqlc.yaml"
			//}

			if queries {
				if err = Queries(QueriesOptions{
					Auto:     options.Auto,
					Sqlc:     options.Sqlc,
					Platform: options.Platform,
					SqlcYaml: yaml,
				}); err != nil {
					return
				}
			}
		}
	}

	err = FixImports(FixImportsOptions{Directory: filepath.Join("lib", "database")})

	return
}
