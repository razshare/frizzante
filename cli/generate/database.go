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
		err = fmt.Errorf("database of type %s is not supported", options.Type)
		return
	}

	dbtype := strings.ToLower(options.Type)
	from := "internal/additions/lib/database/" + dbtype
	to := filepath.Join("lib", "database", dbtype)

	if files.IsDirectory(to) {
		if !options.Auto {
			var yes bool
			if yes, err = confirm.Sendf(true, "%s already exists. Overwrite?", to); err != nil {
				return
			}

			if !yes {
				messages.Infof("skipping %s", to)
				return nil
			}
		}

		if err = os.RemoveAll(to); err != nil {
			return
		}
	}

	if err = Copy(CopyOptions{From: from, To: to, Auto: options.Auto, Efs: options.Efs}); err != nil {
		return
	}

	if dbtype == "sqlite" {
		spin := spinner.New("adding github.com/mattn/go-sqlite3")

		go spinner.Start(spin)
		if !messages.Command(".", os.Environ(), options.Go, "get", "github.com/mattn/go-sqlite3") {
			spinner.Stop(spin)
			err = errors.New("could not add github.com/mattn/go-sqlite3")
			return
		}
		spinner.Stop(spin)

		messages.Success("sqlite database is ready")

		if err = FixImports(FixImportsOptions{Directory: to}); err != nil {
			return
		}

		queries := strings.Contains(strings.ToLower(options.Generate), "queries")

		if !queries {
			if options.Auto {
				queries = false
			} else if queries, err = confirm.Send(true, "would you like to also generate your queries?"); err != nil {
				return
			}
			if queries {
				if err = Queries(QueriesOptions{Auto: options.Auto, Sqlc: options.Sqlc, Platform: options.Platform}); err != nil {
					return
				}
			}
		}
	}

	err = FixImports(FixImportsOptions{Directory: filepath.Join("lib", "database")})

	return
}
