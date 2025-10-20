package database

import (
	"log"

	"github.com/razshare/frizzante/internal/additions/lib/database/sqlite/sqlc"
)

var Queries *sqlc.Queries

func init() {
	if database, err := Load(); err != nil {
		log.Fatal(err)
	} else {
		Queries = sqlc.New(database)
	}
}
