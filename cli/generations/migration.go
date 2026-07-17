package generations

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/razshare/frizzante/internal/project/lib/core/files"
	"github.com/razshare/frizzante/tui/inputs"
	"github.com/razshare/frizzante/tui/messages"
	"github.com/razshare/frizzante/tui/search"
	"github.com/razshare/frizzante/tui/select_one"
	"github.com/razshare/frizzante/tui/spinners"
)

func Migration(options MigrationOptions) (err error) {
	sqlcYaml := options.SqlcYaml
	if sqlcYaml == "" {
		if options.Strict {
			err = errors.New("no sqlc.yaml file provided")
			return
		}
		var names []string
		if names, err = files.FindWithSuffix("lib", "sqlc.yaml"); err != nil {
			return
		}
		choices := make([]search.Choice, len(names))
		for index, name := range names {
			choices[index] = search.Choice{Id: name}
		}
		sqlcYaml, err = select_one.Sendf(choices, "where is your sqlc.yaml file located?")
		if sqlcYaml == "other" {
			if sqlcYaml, err = inputs.Send("where is your sqlc.yaml file located?"); err != nil {
				return
			}
		}
	}
	type Configuration struct {
		Sql []struct {
			Schema string `yaml:"schema"`
		} `yaml:"sql"`
	}
	baseDirectory := filepath.Join("migrate")
	spin := spinners.New("creating migration file")
	go spinners.Start(spin)
	if !files.IsDirectory(filepath.Join(baseDirectory, "migrations")) {
		if err = os.MkdirAll(filepath.Join(baseDirectory, "migrations"), os.ModePerm); err != nil {
			spinners.Stop(spin)
			return
		}
	}
	data := []byte("-- migration: down\n\n-- migration: up\n")
	now := time.Now()
	format := now.Format("2006_01_02T15_04_05Z07_00")
	migrationFileName := filepath.Join(baseDirectory, "migrations", fmt.Sprintf("%s.sql", format))
	if err = os.WriteFile(migrationFileName, data, os.ModePerm); err != nil {
		spinners.Stop(spin)
		return
	}
	spinners.Stop(spin)
	messages.Successf("migration generated at %s", migrationFileName)
	return
}
