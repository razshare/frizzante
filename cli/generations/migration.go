package generations

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/razshare/frizzante/internal/project/lib/core/files"
	"github.com/razshare/frizzante/tui/inputs"
	"github.com/razshare/frizzante/tui/messages"
	"github.com/razshare/frizzante/tui/search"
	"github.com/razshare/frizzante/tui/select_one"
	"github.com/razshare/frizzante/tui/spinners"
	"gopkg.in/yaml.v3"
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
	baseDirectory := filepath.Dir(options.SqlcYaml)
	spin := spinners.New("checking sql code")
	go spinners.Start(spin)
	if !messages.Command(messages.CommandOptions{
		DirectoryName: baseDirectory,
		Environment:   os.Environ(),
		Program:       options.Sqlc,
		Args:          []string{"vet"},
	}) {
		spinners.Stop(spin)
		err = errors.New("sql code check failed")
		return
	}
	spinners.Stop(spin)
	messages.Success("sql code check succeeded")
	var migrationFileName string
	spin = spinners.New("creating migration file")
	go spinners.Start(spin)
	defer func() {
		if err == nil {
			messages.Successf("migration generated at %s", migrationFileName)
		}
	}()
	var data []byte
	if data, err = os.ReadFile(options.SqlcYaml); err != nil {
		spinners.Stop(spin)
		return
	}
	var config Configuration
	if err = yaml.Unmarshal(data, &config); err != nil {
		spinners.Stop(spin)
		return
	}
	if !files.IsDirectory(filepath.Join(baseDirectory, "migrations")) {
		if err = os.MkdirAll(filepath.Join(baseDirectory, "migrations"), os.ModePerm); err != nil {
			spinners.Stop(spin)
			return
		}
	}
	var names []string
	if names, err = files.ReadDirectory(filepath.Join(baseDirectory, "migrations")); err != nil {
		spinners.Stop(spin)
		return
	}
	if len(config.Sql) == 0 {
		spinners.Stop(spin)
		err = errors.New("sql schema not found in configuration file")
		return
	}
	if len(names) == 0 {
		schema := config.Sql[0].Schema
		if !strings.HasSuffix(schema, ".sql") {
			spinners.Stop(spin)
			err = errors.New("database schema up file must have suffix .sql")
			return
		}
		if data, err = os.ReadFile(filepath.Join(baseDirectory, schema)); err != nil {
			spinners.Stop(spin)
			return
		}
	} else {
		data = []byte("-- migration: down\n\n-- migration: up\n")
	}
	now := time.Now()
	migrationFileName = filepath.Join(baseDirectory, "migrations", fmt.Sprintf("%s.sql", now.Format("2006-01-02T15:04:05Z07:00")))
	if err = os.WriteFile(migrationFileName, data, os.ModePerm); err != nil {
		spinners.Stop(spin)
		return
	}
	spinners.Stop(spin)
	return
}
