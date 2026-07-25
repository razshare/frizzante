package main

import (
	"database/sql"
	"embed"
	"log"

	"github.com/razshare/frizzante/v2/internal/project/lib/databases"
)

//go:embed migrations
var efs embed.FS

func main() {
	// no need to close database connection,
	// this program runs once and dies immediately
	var err error
	var database *sql.DB
	if database, err = databases.Connect(); err != nil {
		log.Fatal(err)
	}
	if err = databases.Migrate(databases.MigrateOptions{
		Efs:      efs,
		Database: database,
		Offset:   "first",
		Target:   "last",
	}); err != nil {
		log.Fatal(err)
	}
}
