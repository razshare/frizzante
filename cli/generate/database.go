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
	"github.com/razshare/frizzante/tui/select_one"
	"github.com/razshare/frizzante/tui/spinners"
)

func Database(options DatabaseOptions) (err error) {
	if options.Type == "" {
		if options.Auto {
			options.Type = "sqlite"
		} else if options.Type, err = select_one.Send([]search.Choice{{Id: "sqlite"}}, "what type of database would you like to setup?"); err != nil {
			return
		}
	}

	if options.Type != "sqlite" {
		err = fmt.Errorf("database of type %s is not supported", options.Type)
		return
	}

	dbtype := strings.ToLower(options.Type)
	from := fmt.Sprintf("internal/additions/lib/%s/databases", dbtype)
	to := filepath.Join("lib", dbtype, "databases")

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
		spin := spinners.New("adding github.com/mattn/go-sqlite3")

		go spinners.Start(spin)
		if !messages.Command(messages.CommandOptions{
			Dir:  "app",
			Env:  os.Environ(),
			Name: options.Go,
			Args: []string{"get", "github.com/mattn/go-sqlite3"},
		}) {
			spinners.Stop(spin)
			err = errors.New("could not add github.com/mattn/go-sqlite3")
			return
		}

		spinners.Stop(spin)

		messages.Success("sqlite database is ready")

		if err = FixImports(FixImportsOptions{Directory: to}); err != nil {
			return
		}

		if migration := strings.Contains(strings.ToLower(options.Generate), "migration"); !migration {
			if options.Auto {
				migration = false
			} else if migration, err = confirm.Send(true, "would you like to generate your first migration?"); err != nil {
				return
			}
			if migration {
				if err = Migration(MigrationOptions{Auto: options.Auto, Sqlc: options.Sqlc, Platform: options.Platform}); err != nil {
					return
				}
			}
		}

		if queries := strings.Contains(strings.ToLower(options.Generate), "queries"); !queries {
			if options.Auto {
				queries = false
			} else if queries, err = confirm.Send(true, "would you like to generate your queries?"); err != nil {
				return
			}
			if queries {
				if err = Queries(QueriesOptions{
					Auto:     options.Auto,
					Sqlc:     options.Sqlc,
					Platform: options.Platform,
				}); err != nil {
					return
				}
			}
		}
	}

	err = FixImports(FixImportsOptions{Directory: filepath.Join("lib", "sqlite")})

	return
}
