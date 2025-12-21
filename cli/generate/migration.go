package generate

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/razshare/frizzante/internal/project/lib/core/files"
	"github.com/razshare/frizzante/tui/messages"
	"github.com/razshare/frizzante/tui/spinners"
	"gopkg.in/yaml.v3"
)

func Migration(options MigrationOptions) (err error) {
	type Configuration struct {
		Sql []struct {
			Schema string `yaml:"schema"`
		} `yaml:"sql"`
	}

	baseDirectory := filepath.Dir(options.SqlcYaml)

	if _, err = exec.LookPath(options.Sqlc); err != nil && !files.IsFile(options.Sqlc) {
		if err = Sqlc(SqlcOptions{
			Sqlc:     options.Sqlc,
			Platform: options.Platform,
		}); err != nil {
			return
		}
	}

	var sqlc string
	if files.IsFile(options.Sqlc) {
		if sqlc, err = filepath.Rel(baseDirectory, options.Sqlc); err != nil {
			return err
		}
	} else if sqlc, err = exec.LookPath(options.Sqlc); err != nil {
		sqlc = options.Sqlc
	}

	spin := spinners.New("checking sql code")

	go spinners.Start(spin)
	if !messages.Command(messages.CommandOptions{
		DirectoryName: baseDirectory,
		Environment:   os.Environ(),
		Program:       sqlc,
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
