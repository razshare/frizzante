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
	"github.com/razshare/frizzante/tui/messages"
	"github.com/razshare/frizzante/tui/search"
	"github.com/razshare/frizzante/tui/select_one"
	"github.com/razshare/frizzante/tui/spinners"
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

		if !options.Interactive {
			yamlFileName = choices[0].Id
		} else {
			yamlFileName, err = select_one.Sendf(choices, "where is your sqlc.yaml file located?")
		}
	}

	baseDirectory := filepath.Dir(yamlFileName)

	if !files.IsFile(yamlFileName) {
		return fmt.Errorf("%s not found", yamlFileName)
	}

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
		if times[index], err = time.Parse("2006-01-02T15:04:05Z07:00", strings.TrimSuffix(filepath.Base(name), ".sql")); err != nil {
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
	offset := options.Offset
	if offset == "" {
		choices := make([]search.Choice, count)

		for index, time_ := range times {
			choices[index] = search.Choice{Id: time_.Format("2006-01-02T15:04:05Z07:00")}
		}

		if offset, err = select_one.Sendf(choices, "what's the offset migration? (sorting desc)"); err != nil {
			return err
		}
		if offsetTime, err = time.Parse("2006-01-02T15:04:05Z07:00", offset); err != nil {
			return
		}
	} else if offset == "first" {
		offsetTime = times[count-1]
	} else if offset == "last" {
		offsetTime = times[0]
	} else if offsetTime, err = time.Parse("2006-01-02T15:04:05Z07:00", offset); err != nil {
		return
	}

	var targetTime time.Time
	target := options.Target
	if target == "" {
		choices := make([]search.Choice, count)

		for index, time_ := range times {
			choices[index] = search.Choice{Id: time_.Format("2006-01-02T15:04:05Z07:00")}
		}

		if target, err = select_one.Sendf(choices, "what's the target migration? (sorting desc)"); err != nil {
			return err
		}
		if targetTime, err = time.Parse("2006-01-02T15:04:05Z07:00", target); err != nil {
			return
		}
	} else if target == "first" {
		targetTime = times[count-1]
	} else if target == "last" {
		targetTime = times[0]
	} else if targetTime, err = time.Parse("2006-01-02T15:04:05Z07:00", target); err != nil {
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
			migrations = append(migrations, filepath.Join(baseDirectory, "migrations", value.Format("2006-01-02T15:04:05Z07:00")+".sql"))
		} else {
			if value.Before(targetTime) || value.After(offsetTime) {
				continue
			}
			migrations = append(migrations, filepath.Join(baseDirectory, "migrations", value.Format("2006-01-02T15:04:05Z07:00")+".sql"))
		}
	}

	if len(migrations) == 0 {
		messages.Info("no migrations matched for execution")
		return
	}

	var transaction *sql.Tx
	if transaction, err = options.Database.Begin(); err != nil {
		return
	}

	for _, migration := range migrations {
		spin := spinners.Newf("migrating database schema using %s", migration)
		go spinners.Start(spin)

		var data []byte
		if data, err = os.ReadFile(migration); err != nil {
			spinners.Stop(spin)
			return
		}

		if forward {
			messages.Infof("running %s up", migration)
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

			if query := builder.String(); query != "" {
				if _, err = transaction.Exec(query); err != nil {
					if err = transaction.Rollback(); err != nil {
						spinners.Stop(spin)
						return
					}
					spinners.Stop(spin)
					return
				}
			}
		} else {
			messages.Infof("running %s down", migration)
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

			if query := builder.String(); query != "" {
				if _, err = transaction.Exec(query); err != nil {
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

	messages.Success("database schema migrated successfully")

	return
}
