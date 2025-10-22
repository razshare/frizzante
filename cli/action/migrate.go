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

	"github.com/razshare/frizzante/cli/generate"
	"github.com/razshare/frizzante/internal/project/lib/core/files"
	"github.com/razshare/frizzante/tui/messages"
	"github.com/razshare/frizzante/tui/multiselect"
	"github.com/razshare/frizzante/tui/search"
	"github.com/razshare/frizzante/tui/singleselect"
	"github.com/razshare/frizzante/tui/spinner"
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

		choices = append(choices, search.Choice{Id: "other", Description: "use a different file"})

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
	if names, err = files.ReadDirectory(filepath.Join(baseDirectory, "migrations")); err != nil {
		return
	}

	count := len(names)
	if count == 0 {
		err = errors.New("no migration files found")
		return
	}

	mixed := map[int]string{}
	for _, name := range names {
		trimmed := strings.TrimSuffix(filepath.Base(name), ".sql")
		parts := strings.SplitN(trimmed, "_", 2)
		if len(parts) < 1 {
			err = fmt.Errorf("invalid migration file name %s", name)
			return
		}

		var parsed int64
		if parsed, err = strconv.ParseInt(parts[0], 10, 64); err != nil {
			return
		}

		key := int(parsed)

		if _, exists := mixed[key]; exists {
			err = fmt.Errorf("duplicate migration key %d", key)
			return
		}

		mixed[key] = name
	}

	keys := make([]int, 0)
	for id := range mixed {
		keys = append(keys, id)
	}

	slices.Sort(keys)

	sorted := map[int]string{}
	for index, key := range keys {
		sorted[index] = mixed[key]
	}

	queryString := options.QueryString
	migrations := make([]string, 0)

	if queryString == "" {
		choices := []search.Choice{
			{Id: "latest", Description: "only latest"},
			{Id: "all", Description: "in order from first to last"},
			{Id: "after", Description: "in order from offest (exclusive) to latest (inclusive)"},
			{Id: "before", Description: "in order from first (inclusive) to offset (exclusive)"},
		}

		if queryString, err = singleselect.Sendf(choices, "which migrations would you like to execute?"); err != nil {
			return err
		}

		if queryString == "after" {
			queryString = ">"
		} else if queryString == "before" {
			queryString = "<"
		}
	}

	if queryString == "pick" {
		choices := make([]search.Choice, 0)
		for _, name := range sorted {
			choices = append(choices, search.Choice{Id: name})
		}
		slices.Reverse(choices)
		if migrations, err = multiselect.Sendf(choices, "select a migration to execute"); err != nil {
			return err
		}
	} else if queryString == "*" || queryString == "all" {
		for _, name := range sorted {
			migrations = append(migrations, name)
		}
	} else if strings.HasPrefix(queryString, ">") {
		if queryString == ">" {
			choices := make([]search.Choice, len(sorted))
			for key, name := range sorted {
				choices[key] = search.Choice{Id: name}
			}
			slices.Reverse(choices)
			if queryString, err = singleselect.Send(choices, "pick an offset (exclusive)"); err != nil {
				return err
			}
			queryString = ">" + strings.SplitN(filepath.Base(queryString), "_", 2)[0]
		}

		var value int64
		if value, err = strconv.ParseInt(queryString[1:], 10, 64); err != nil {
			return
		}

		for key, name := range sorted {
			if key > int(value-1) {
				migrations = append(migrations, name)
			}
		}
	} else if strings.HasPrefix(queryString, "<") {
		if queryString == "<" {
			choices := make([]search.Choice, len(sorted))
			for key, name := range sorted {
				choices[key] = search.Choice{Id: name}
			}
			slices.Reverse(choices)
			if queryString, err = singleselect.Send(choices, "pick an offset (exclusive)"); err != nil {
				return err
			}
			queryString = "<" + strings.SplitN(filepath.Base(queryString), "_", 2)[0]
		}

		var value int64
		if value, err = strconv.ParseInt(queryString[1:], 10, 64); err != nil {
			return
		}

		for key, name := range sorted {
			if key < int(value-1) {
				migrations = append(migrations, name)
			}
		}
	} else if queryString == "latest" {
		var value int
		for key := range sorted {
			if key > value {
				value = key
			}
		}
		migrations = append(migrations, sorted[value])
	} else {
		err = errors.New("unknown migrate query string")
		return
	}

	var transaction *sql.Tx
	if transaction, err = options.Database.Begin(); err != nil {
		return
	}

	if len(migrations) == 0 {
		messages.Info("no migrations matched for execution")
		return
	}

	for _, migration := range migrations {
		spin := spinner.Newf("migrating database schema using %s", migration)
		go spinner.Start(spin)

		var data []byte
		if data, err = os.ReadFile(migration); err != nil {
			spinner.Stop(spin)
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

		if query := down.String(); query != "" {
			if _, err = transaction.Exec(query); err != nil {
				if err = transaction.Rollback(); err != nil {
					spinner.Stop(spin)
					return
				}
				spinner.Stop(spin)
				return
			}
		}

		if query := up.String(); query != "" {
			if _, err = transaction.Exec(query); err != nil {
				if err = transaction.Rollback(); err != nil {
					spinner.Stop(spin)
					return
				}
				spinner.Stop(spin)
				return
			}
		}
		spinner.Stop(spin)
		messages.Successf("migration %s executed successfully", migration)
	}

	if err = transaction.Commit(); err != nil {
		return
	}

	messages.Success("database schema migrated successfully")

	return
}
