package generations

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/razshare/frizzante/tui/messages"
	"github.com/razshare/frizzante/tui/search"
	"github.com/razshare/frizzante/tui/select_one"
	"github.com/razshare/frizzante/tui/spinners"
)

func Databases(options DatabasesOptions) (err error) {
	databaseType := strings.ToLower(options.Type)
	if databaseType == "" {
		if options.Strict {
			err = errors.New("no database type provided")
			return
		}
		if databaseType, err = select_one.Send(
			[]search.Choice{{Id: "sqlite"}},
			"what type of database would you like to setup?",
		); err != nil {
			return
		}
	}
	fromDirectoryName := fmt.Sprintf("internal/additions/lib/databases/%s", databaseType)
	toDirectoryName := filepath.Join("lib", "databases", databaseType)
	if err = Copy(CopyOptions{
		From: fromDirectoryName,
		To:   toDirectoryName,
		Efs:  options.Efs,
	}); err != nil {
		return
	}
	if databaseType != "sqlite" {
		err = fmt.Errorf("%s database type is not supported", databaseType)
		return
	}
	spin := spinners.New("adding github.com/mattn/go-sqlite3")
	go spinners.Start(spin)
	if !messages.Command(messages.CommandOptions{
		Environment: os.Environ(),
		Program:     options.Go,
		Args:        []string{"get", "github.com/mattn/go-sqlite3"},
	}) {
		spinners.Stop(spin)
		err = errors.New("could not add github.com/mattn/go-sqlite3")
		return
	}
	spinners.Stop(spin)
	messages.Success("sqlite database is ready")
	if err = Queries(QueriesOptions{
		Sqlc:     options.Sqlc,
		SqlcYaml: filepath.Join(toDirectoryName, "sqlc.yaml"),
		Strict:   options.Strict,
	}); err != nil {
		return
	}
	err = FixImports(FixImportsOptions{Directory: toDirectoryName})
	return
}
