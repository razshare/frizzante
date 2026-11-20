package databases

import (
	"errors"
	"fmt"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/razshare/frizzante/v2/internal/project/lib/core/embeds"
)

func ListMigrations(options ListMigrationsOptions) (err error) {
	efs := options.Efs
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
	red := "\x1b[31m"
	green := "\x1b[32m"
	reset := "\x1b[0m"
	for _, version := range migrations {
		fmt.Printf("version %s...", version)
		var exists bool
		if exists, err = options.Exists(version); err != nil {
			return
		}
		if exists {
			fmt.Printf(" %sfound%s.\n", green, reset)
		} else {
			fmt.Printf(" %smissing%s.\n", red, reset)
		}
	}
	return
}
