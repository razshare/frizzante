package generate

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/razshare/frizzante/internal/project/lib/core/files"
	"github.com/razshare/frizzante/tui/messages"
	"github.com/razshare/frizzante/tui/search"
	"github.com/razshare/frizzante/tui/singleselect"
	"github.com/razshare/frizzante/tui/spinner"
	"gopkg.in/yaml.v3"
)

func Migration(options MigrationOptions) (err error) {
	type Configuration struct {
		Sql []struct {
			Schema string `yaml:"schema"`
		} `yaml:"sql"`
	}

	yamlFileName := options.SqlcYaml

	if yamlFileName != "" && !files.IsFile(yamlFileName) {
		messages.Infof("%s not found", yamlFileName)
	}

	if yamlFileName == "" {
		var items []string
		if items, err = files.ReadDirectory("lib"); err != nil {
			return
		}

		names := make([]string, 0)
		for _, item := range items {
			if strings.HasSuffix(item, string(filepath.Separator)+"sqlc.yaml") {
				names = append(names, item)
			}
		}

		choices := make([]search.Choice, len(names))
		for index, name := range names {
			choices[index] = search.Choice{Id: name}
		}

		choices = append(choices, search.Choice{Id: "other", Description: "other"})

		if options.Auto {
			yamlFileName = choices[0].Id
		} else {
			yamlFileName, err = singleselect.Sendf(choices, "where is your sqlc.yaml file located?")
		}
	}

	baseDirectory := filepath.Dir(yamlFileName)

	if _, err = exec.LookPath(options.Sqlc); err != nil && !files.IsFile(options.Sqlc) {
		if err = Sqlc(SqlcOptions{
			Sqlc:     options.Sqlc,
			Platform: options.Platform,
			Auto:     options.Auto,
		}); err != nil {
			return
		}
	}

	if !files.IsFile(yamlFileName) {
		return fmt.Errorf("%s not found", yamlFileName)
	}

	var sqlc string
	if files.IsFile(options.Sqlc) {
		if sqlc, err = filepath.Rel(baseDirectory, options.Sqlc); err != nil {
			return err
		}
	} else if sqlc, err = exec.LookPath(options.Sqlc); err != nil {
		sqlc = options.Sqlc
	}

	spin := spinner.New("checking sql code")

	go spinner.Start(spin)
	if !messages.Command(baseDirectory, os.Environ(), sqlc, "vet") {
		spinner.Stop(spin)
		err = errors.New("sql code check failed")
		return
	}
	spinner.Stop(spin)
	messages.Success("sql code check succeeded")

	var migrationFileName string

	spin = spinner.New("creating migration file")
	go spinner.Start(spin)
	defer func() {
		if err == nil {
			messages.Successf("migration generated at %s", migrationFileName)
		}
	}()
	defer spinner.Stop(spin)

	var data []byte
	if data, err = os.ReadFile(yamlFileName); err != nil {
		return
	}
	var config Configuration
	if err = yaml.Unmarshal(data, &config); err != nil {
		return
	}

	var names []string
	if names, err = files.ReadDirectory(filepath.Join(baseDirectory, "migrations")); err != nil {
		return
	}

	numbers := make([]int64, len(names))

	var index int64
	for jndex, entry := range names {
		trimmed := strings.TrimSuffix(filepath.Base(entry), ".sql")
		parts := strings.SplitN(trimmed, "_", 2)
		if len(parts) < 1 {
			err = errors.New("invalid empty migration name")
			return
		}

		var value int64
		if value, err = strconv.ParseInt(parts[0], 10, 64); err != nil {
			return
		}

		if slices.Contains(numbers, value) {
			err = fmt.Errorf("duplicate migration index %s", parts[0])
			return
		}

		numbers[jndex] = value

		if index < value {
			index = value
		}
	}

	if len(config.Sql) == 0 {
		err = errors.New("sql schema not found in configuration file")
		return
	}

	schemaUpFileName := config.Sql[0].Schema
	if !strings.HasSuffix(schemaUpFileName, ".up.sql") {
		err = errors.New("database schema up file must have suffix .up.sql")
		return
	}

	schemaDownFileName := strings.TrimSuffix(config.Sql[0].Schema, ".up.sql") + ".down.sql"

	var migration strings.Builder
	if data, err = os.ReadFile(filepath.Join(baseDirectory, schemaDownFileName)); err != nil {
		return
	}
	migration.WriteString("-- migration: down\n")
	migration.WriteString(strings.TrimSpace(string(data)) + "\n\n")

	if data, err = os.ReadFile(filepath.Join(baseDirectory, schemaUpFileName)); err != nil {
		return
	}
	migration.WriteString("-- migration: up\n")
	migration.WriteString(strings.TrimSpace(string(data)) + "\n")

	now := time.Now()
	migrationFileName = filepath.Join(baseDirectory, "migrations", fmt.Sprintf("%d_%s.sql", index+1, now.Format("2006-01-02 15.04.05.00")))

	if err = os.WriteFile(migrationFileName, []byte(migration.String()), os.ModePerm); err != nil {
		return
	}

	return
}
