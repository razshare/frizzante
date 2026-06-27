package main

import (
	"database/sql"
	"log"

	"github.com/razshare/frizzante/internal/project/lib/core/databases"
)

func main() {
	// no need to close database connection,
	// this program runs once and dies immediately
	var err error
	var database *sql.DB
	if database, _, err = databases.Connect(); err != nil {
		log.Fatal(err)
	}
	if err = databases.Migrate(database, "first", "last"); err != nil {
		log.Fatal(err)
	}
}
