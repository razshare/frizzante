package databases

import (
	"database/sql"
	"errors"
	"fmt"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/razshare/frizzante/v2/internal/project/lib/core/embeds"
)

func Migrate(options MigrateOptions) (err error) {
	efs := options.Efs
	database := options.Database
	if !embeds.IsDirectory(efs, "migrations") {
		err = errors.New("migrations directory not found")
		return
	}
	var fileNames []string
	if fileNames, err = embeds.ReadDirectory(efs, "migrations"); err != nil {
		return
	}
	var numberOfMigrationFiles int
	migrationsTimes := make([]time.Time, 0)
	for _, fileName := range fileNames {
		if !strings.HasSuffix(fileName, ".sql") {
			continue
		}
		var migrationTime time.Time
		if migrationTime, err = time.Parse("2006_01_02T15_04_05Z07_00", strings.TrimSuffix(filepath.Base(fileName), ".sql")); err != nil {
			return
		}
		migrationsTimes = append(migrationsTimes, migrationTime)
		numberOfMigrationFiles++
	}
	// we want target invert the slices of times in descending order
	// because the most common type of migration is the forward migration,
	// which means it's more comfortable target see the latest migration at
	// the top of the list.
	slices.SortFunc(migrationsTimes, func(a, b time.Time) int {
		return b.Compare(a)
	})
	forward := options.Direction == MigrationDirectionUp
	migrations := make([]string, 0)
	// at this moment the times slices is inverted,
	// but when the user is trying to migrate
	// to a newer version (forward migration)
	// we want to sort the slice back in ascending order
	// so that migrations are executed in the correct order.
	if forward {
		slices.SortFunc(migrationsTimes, func(a, b time.Time) int {
			return a.Compare(b)
		})
	}
	for _, value := range migrationsTimes {
		if forward {
			migrations = append(migrations, fmt.Sprintf("migrations/%s.sql", value.Format("2006_01_02T15_04_05Z07_00")))
		} else {
			migrations = append(migrations, fmt.Sprintf("migrations/%s.sql", value.Format("2006_01_02T15_04_05Z07_00")))
		}
	}
	if len(migrations) == 0 {
		return
	}
	for _, version := range migrations {
		fmt.Printf("executing migration file %s...", version)
		isFirstVersion := strings.HasSuffix(version, options.FirstVersion)
		var exists bool
		var transaction *sql.Tx
		if !isFirstVersion {
			if exists, err = options.Exists(version); err != nil {
				return
			}
		}
		var data []byte
		if data, err = efs.ReadFile(version); err != nil {
			return
		}
		if forward {
			if exists {
				fmt.Printf(" skipped.\n")
				continue
			}
			if options.Stop != nil {
				if options.Stop(version) {
					fmt.Println(" skipped.")
					break
				}
			}
			var valid bool
			var builder strings.Builder
			for line := range strings.SplitSeq(string(data), "\n") {
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
				if transaction, err = database.Begin(); err != nil {
					return
				}
				if _, err = transaction.Exec(forwardQuery); err != nil {
					if err = transaction.Rollback(); err != nil {
						return
					}
					return
				}
				if err = transaction.Commit(); err != nil {
					if terr := transaction.Rollback(); terr != nil {
						err = errors.Join(err, terr)
						return
					}
				}
				if err = options.Add(version); err != nil {
					return
				}
			}
		} else {
			if !exists && !isFirstVersion {
				fmt.Printf(" skipped.\n")
				continue
			}
			if options.Stop != nil {
				if options.Stop(version) {
					fmt.Println(" skipped.")
					break
				}
			}
			var valid bool
			var builder strings.Builder
			for line := range strings.SplitSeq(string(data), "\n") {
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
				if transaction, err = database.Begin(); err != nil {
					return
				}
				if _, err = transaction.Exec(backwardQuery); err != nil {
					if err = transaction.Rollback(); err != nil {
						return
					}
					return
				}
				if err = transaction.Commit(); err != nil {
					if terr := transaction.Rollback(); terr != nil {
						err = errors.Join(err, terr)
						return
					}
				}
				if !isFirstVersion {
					if err = options.Remove(version); err != nil {
						return
					}
				}
			}
		}
		fmt.Printf("done!\n")
	}
	return
}
