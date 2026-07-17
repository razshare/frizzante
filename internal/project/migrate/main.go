package main

import (
	"database/sql"
	"embed"
	"log"

	databases2 "github.com/razshare/frizzante/internal/project/lib/databases"
)

//go:embed migrations
var efs embed.FS

func main() {
	// no need to close database connection,
	// this program runs once and dies immediately
	var err error
	var database *sql.DB
	if database, _, err = databases2.Connect(); err != nil {
		log.Fatal(err)
	}
	if err = databases2.Migrate(databases2.MigrateOptions{
		Efs:      efs,
		Database: database,
		Offset:   "first",
		Target:   "last",
	}); err != nil {
		log.Fatal(err)
	}
}
