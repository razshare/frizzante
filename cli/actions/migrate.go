package actions

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/razshare/frizzante/internal/project/lib/core/files"
	"github.com/razshare/frizzante/tui/inputs"
	"github.com/razshare/frizzante/tui/messages"
	"github.com/razshare/frizzante/tui/search"
	"github.com/razshare/frizzante/tui/select_many"
	"github.com/razshare/frizzante/tui/select_one"
	"github.com/razshare/frizzante/tui/spinners"
)

// Migrate runs the latest migration file against the given database.
//
// Experimental: api is currently minimal and not stable.
func Migrate(options MigrateOptions) (err error) {
	databaseConnectionStrings := strings.Split(options.Database, ",")
	if len(databaseConnectionStrings) == 1 && strings.TrimSpace(databaseConnectionStrings[0]) == "" {
		databaseConnectionStrings = make([]string, 0)
	}
	sqlcYamlFileName := options.SqlcYaml
	query := options.Query
	if sqlcYamlFileName == "" {
		if options.Strict {
			err = errors.New("no sqlc.yaml file provided")
			return
		}
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
		sqlcYamlFileName, err = select_one.Sendf(choices, "where is your sqlc.yaml file located?")
		if sqlcYamlFileName == "other" {
			if sqlcYamlFileName, err = inputs.Send("where is your sqlc.yaml file located?"); err != nil {
				return
			}
		}
	}
	var offset string
	var target string
	migrateRange := strings.SplitN(query, ",", 2)
	if len(migrateRange) >= 1 {
		offset = migrateRange[0]
	} else {
		offset = ""
	}
	if len(migrateRange) >= 2 {
		target = migrateRange[1]
		if offset == "" {
			offset = "first"
		}
		if target == "" {
			target = "last"
		}
	} else {
		target = ""
	}
	if len(databaseConnectionStrings) == 0 {
		if options.Strict {
			err = errors.New("database connection string not provided")
			return
		}
		var databaseFileNames []string
		if databaseFileNames, err = files.FindWithSuffix(".", ".sqlite"); err != nil {
			return
		}
		choices := make([]search.Choice, len(databaseFileNames))
		for index, databaseFileName := range databaseFileNames {
			choices[index] = search.Choice{Id: databaseFileName}
		}
		var selectedDatabaseFileNames []string
		for {
			if selectedDatabaseFileNames, err = select_many.Sendf(choices, "pick databases to migrate"); err != nil {
				return
			}
			if len(selectedDatabaseFileNames) == 0 {
				messages.Info("you must pick at least one database to migrate")
				continue
			}
			for _, selectedDatabaseFilename := range selectedDatabaseFileNames {
				databaseConnectionStrings = append(databaseConnectionStrings, fmt.Sprintf("file:%s?cache=shared", selectedDatabaseFilename))
			}
			break
		}
	}
	migrate := func(databaseConnectionString string) (err error) {
		var databaseConnection *sql.DB
		if databaseConnection, err = sql.Open("sqlite3", databaseConnectionString); err != nil {
			return
		}
		defer func() {
			if cerr := databaseConnection.Close(); cerr != nil {
				if err == nil {
					err = cerr
				}
			}
		}()
		baseDirectory := filepath.Dir(sqlcYamlFileName)
		if !files.IsDirectory(filepath.Join(baseDirectory, "migrations")) {
			if err = os.MkdirAll(filepath.Join(baseDirectory, "migrations"), os.ModePerm); err != nil {
				return
			}
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
		times := make([]time.Time, count)
		for index, name := range names {
			if times[index], err = time.Parse("2006_01_02T15_04_05Z07_00", strings.TrimSuffix(filepath.Base(name), ".sql")); err != nil {
				return
			}
		}
		// we want to invert the slices of times in descending order
		// because the most common type of migration is the forward migration,
		// which means it's more comfortable to see the latest migration at
		// the top of the list.
		slices.SortFunc(times, func(a, b time.Time) int {
			return b.Compare(a)
		})
		var offsetTime time.Time
		if offset == "" {
			choices := make([]search.Choice, count)
			for index, time_ := range times {
				choices[index] = search.Choice{Id: time_.Format("2006_01_02T15_04_05Z07_00")}
			}
			if offset, err = select_one.Sendf(choices, "what's the offset migration? (sorting desc)"); err != nil {
				return
			}
			if offsetTime, err = time.Parse("2006_01_02T15_04_05Z07_00", offset); err != nil {
				return
			}
		} else if offset == "first" {
			offsetTime = times[count-1]
		} else if offset == "last" {
			offsetTime = times[0]
		} else if offsetTime, err = time.Parse("2006_01_02T15_04_05Z07_00", offset); err != nil {
			return
		}
		var targetTime time.Time
		if target == "" {
			choices := make([]search.Choice, count)
			for index, time_ := range times {
				choices[index] = search.Choice{Id: time_.Format("2006_01_02T15_04_05Z07_00")}
			}
			if target, err = select_one.Sendf(choices, "what's the target migration? (sorting desc)"); err != nil {
				return
			}
			if targetTime, err = time.Parse("2006_01_02T15_04_05Z07_00", target); err != nil {
				return
			}
		} else if target == "first" {
			targetTime = times[count-1]
		} else if target == "last" {
			targetTime = times[0]
		} else if targetTime, err = time.Parse("2006_01_02T15_04_05Z07_00", target); err != nil {
			return
		}
		forward := !offsetTime.After(targetTime)
		migrations := make([]string, 0)
		// at this moment the times slices is inverted,
		// but when the user is trying to migrate
		// to a newer version (forward migration)
		// we want to sort the slice back in ascending order
		// so that migrations are executed in the correct order.
		if forward {
			slices.SortFunc(times, func(a, b time.Time) int {
				return a.Compare(b)
			})
		}
		for _, value := range times {
			if forward {
				if value.Before(offsetTime) || value.After(targetTime) {
					continue
				}
				migrations = append(migrations, filepath.Join(baseDirectory, "migrations", value.Format("2006_01_02T15_04_05Z07_00")+".sql"))
			} else {
				if value.Before(targetTime) || value.After(offsetTime) {
					continue
				}
				migrations = append(migrations, filepath.Join(baseDirectory, "migrations", value.Format("2006_01_02T15_04_05Z07_00")+".sql"))
			}
		}
		if len(migrations) == 0 {
			messages.Info("no migrations matched for execution")
			return
		}
		var transaction *sql.Tx
		if transaction, err = databaseConnection.Begin(); err != nil {
			return
		}
		for _, migration := range migrations {
			spin := spinners.Newf("migrating database %s\nusing %s", databaseConnectionString, migration)
			go spinners.Start(spin)
			var data []byte
			if data, err = os.ReadFile(migration); err != nil {
				spinners.Stop(spin)
				return
			}
			if forward {
				var valid bool
				var builder strings.Builder
				for _, line := range strings.Split(string(data), "\n") {
					if line == "-- migration: up" {
						valid = true
						continue
					} else if line == "-- migration: down" {
						valid = false
						continue
					}
					if !valid {
						continue
					}
					builder.WriteString(line)
					builder.WriteString("\n")
				}
				if forwardQuery := builder.String(); forwardQuery != "" {
					if _, err = transaction.Exec(forwardQuery); err != nil {
						if err = transaction.Rollback(); err != nil {
							spinners.Stop(spin)
							return
						}
						spinners.Stop(spin)
						return
					}
				}
			} else {
				var valid bool
				var builder strings.Builder
				for _, line := range strings.Split(string(data), "\n") {
					if line == "-- migration: down" {
						valid = true
						continue
					} else if line == "-- migration: up" {
						valid = false
						continue
					}
					if !valid {
						continue
					}
					builder.WriteString(line)
					builder.WriteString("\n")
				}
				if backwardQuery := builder.String(); backwardQuery != "" {
					if _, err = transaction.Exec(backwardQuery); err != nil {
						if err = transaction.Rollback(); err != nil {
							spinners.Stop(spin)
							return
						}
						spinners.Stop(spin)
						return
					}
				}
			}
			spinners.Stop(spin)
		}
		if err = transaction.Commit(); err != nil {
			return
		}
		return
	}
	for _, databaseConnectionString := range databaseConnectionStrings {
		if err = migrate(databaseConnectionString); err != nil {
			return
		}
	}
	messages.Success("database schema migrated successfully")
	return
}
