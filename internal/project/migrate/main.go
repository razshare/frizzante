package main

import (
	"context"
	"database/sql"
	"embed"
	"flag"
	"log"

	"github.com/razshare/frizzante/v2/internal/project/lib/databases"
	"github.com/razshare/frizzante/v2/internal/project/lib/databases/sqlite"
	"github.com/razshare/frizzante/v2/internal/project/lib/schema"
)

//go:embed migrations
var efs embed.FS
var direction = databases.MigrationDirectionUp

func main() {
	flag.Parse()
	args := flag.Args()
	list := len(args) > 0 && args[0] == "list"
	down := len(args) > 0 && args[0] == "down"
	if down {
		direction = databases.MigrationDirectionDown
	}
	// no need to close database connection,
	// this program runs once and dies immediately
	var err error
	var database *sql.DB
	if database, err = sqlite.Connect(); err != nil {
		log.Fatal(err)
	}
	queries := schema.New(database)
	switch true {
	case list:
		if err = databases.ListMigrations(databases.ListMigrationsOptions{
			Efs:       efs,
			Direction: direction,
			Exists: func(version string) (exists bool, err error) {
				var migration []schema.Migration
				if migration, err = queries.FindMigration(context.Background(), version); err != nil {
					return false, err
				}
				exists = len(migration) > 0
				return
			},
		}); err != nil {
			log.Fatal(err)
		}
	default:
		var counter int
		if err = databases.Migrate(databases.MigrateOptions{
			Efs:          efs,
			Database:     database,
			Direction:    direction,
			FirstVersion: "2026_07_13T09_19_15+02_00.sql",
			Stop: func(_ string) bool {
				counter++
				if direction == databases.MigrationDirectionDown {
					return counter > 1
				}
				return false
			},
			Add: func(version string) (err error) {
				err = queries.AddMigration(context.Background(), version)
				return
			},
			Remove: func(version string) (err error) {
				err = queries.RemoveMigration(context.Background(), version)
				return
			},
			Exists: func(version string) (exists bool, err error) {
				var migration []schema.Migration
				if migration, err = queries.FindMigration(context.Background(), version); err != nil {
					return false, err
				}
				exists = len(migration) > 0
				return
			},
		}); err != nil {
			log.Fatal(err)
		}
	}
}
