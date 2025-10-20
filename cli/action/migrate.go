package action

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strconv"
	"strings"

	_ "github.com/mattn/go-sqlite3"
	"github.com/razshare/frizzante/cli/generate"
	"github.com/razshare/frizzante/internal/project/lib/core/files"
	"github.com/razshare/frizzante/tui/messages"
	"github.com/razshare/frizzante/tui/search"
	"github.com/razshare/frizzante/tui/singleselect"
)

// Migrate runs the latest migration file against the given database.
//
// Experimental: api is currently minimal and not stable.
func Migrate(options MigrateOptions) (err error) {
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
		if err = generate.Sqlc(generate.SqlcOptions{
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

	var names []string
	if names, err = files.ReadDirectory(filepath.Join("lib", "database", "sqlite", "migrations")); err != nil {
		return
	}

	var migrationFileName string

	index := options.Index

	if index == 0 {
		numbers := make([]int64, len(names))
		for jndex, name := range names {
			trimmed := strings.TrimSuffix(filepath.Base(name), ".sql")
			parts := strings.SplitN(trimmed, "_", 2)
			if len(parts) < 1 {
				err = fmt.Errorf("invalid migration file name %s", name)
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
				migrationFileName = name
			}
		}
	} else {
		for _, name := range names {
			if strings.HasPrefix(name, fmt.Sprintf("%d_", index)) {
				migrationFileName = name
				break
			}
		}
	}

	if migrationFileName == "" {
		err = errors.New("migration file name resolved to an empty string")
		return
	}

	messages.Infof("migratind database schema using %s", migrationFileName)

	var data []byte
	if data, err = os.ReadFile(migrationFileName); err != nil {
		return
	}

	var up strings.Builder
	var down strings.Builder
	var tearingDown bool
	for _, line := range strings.Split(string(data), "\n") {
		if line == "-- migration: down" {
			tearingDown = true
			continue
		} else if line == "-- migration: up" {
			tearingDown = false
			continue
		}

		if tearingDown {
			down.WriteString(line)
			down.WriteString("\n")
			continue
		}

		up.WriteString(line)
		up.WriteString("\n")
	}

	var database *sql.DB
	if database, err = sql.Open("sqlite3", fmt.Sprintf("file:%s/source.sqlite?cache=shared", baseDirectory)); err != nil {
		return
	}

	var transaction *sql.Tx
	if transaction, err = database.Begin(); err != nil {
		return
	}

	if query := down.String(); query != "" {
		if _, err = transaction.Exec(query); err != nil {
			if err = transaction.Rollback(); err != nil {
				return
			}
			return
		}
	}

	if query := up.String(); query != "" {
		if _, err = transaction.Exec(query); err != nil {
			if err = transaction.Rollback(); err != nil {
				return
			}
			return
		}
	}

	if err = transaction.Commit(); err != nil {
		return
	}

	messages.Success("migration successful")

	return
}
