package databases

import (
	"database/sql"
	"errors"
	"fmt"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/razshare/frizzante/internal/project/lib/core/embeds"
)

func Migrate(database *sql.DB, offset string, target string) (err error) {
	if offset == "" {
		err = errors.New("offset migration cannot be empty")
		return
	}
	if target == "" {
		err = errors.New("target migration cannot be empty")
		return
	}
	var fileNames []string
	if fileNames, err = embeds.ReadDirectory(Efs, "migrations"); err != nil {
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
	var offsetTime time.Time
	if offset == "first" {
		offsetTime = migrationsTimes[numberOfMigrationFiles-1]
	} else if offset == "last" {
		offsetTime = migrationsTimes[0]
	} else if offsetTime, err = time.Parse("2006_01_02T15_04_05Z07_00", offset); err != nil {
		return
	}
	var targetTime time.Time
	if target == "first" {
		targetTime = migrationsTimes[numberOfMigrationFiles-1]
	} else if target == "last" {
		targetTime = migrationsTimes[0]
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
		slices.SortFunc(migrationsTimes, func(a, b time.Time) int {
			return a.Compare(b)
		})
	}
	for _, value := range migrationsTimes {
		if forward {
			if value.Before(offsetTime) || value.After(targetTime) {
				continue
			}
			migrations = append(migrations, fmt.Sprintf("migrations/%s.sql", value.Format("2006_01_02T15_04_05Z07_00")))
		} else {
			if value.Before(targetTime) || value.After(offsetTime) {
				continue
			}
			migrations = append(migrations, fmt.Sprintf("migrations/%s.sql", value.Format("2006_01_02T15_04_05Z07_00")))
		}
	}
	if len(migrations) == 0 {
		return
	}
	var transaction *sql.Tx
	if transaction, err = database.Begin(); err != nil {
		return
	}

	for _, migration := range migrations {
		fmt.Printf("executing migration file %s... ", migration)
		var data []byte
		if data, err = Efs.ReadFile(migration); err != nil {
			return
		}
		if forward {
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
				if _, err = transaction.Exec(forwardQuery); err != nil {
					if err = transaction.Rollback(); err != nil {
						return
					}
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
						return
					}
					return
				}
			}
		}
		fmt.Printf("done!\n")
	}
	if err = transaction.Commit(); err != nil {
		if terr := transaction.Rollback(); terr != nil {
			err = errors.Join(err, terr)
			return
		}
	}
	return
}
