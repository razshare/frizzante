package databases

import (
	"log"

	"github.com/razshare/frizzante/internal/additions/lib/databases/sqlite/sqlc"
)

var Queries *sqlc.Queries

func init() {
	if database, err := Load(); err != nil {
		log.Fatal(err)
	} else {
		Queries = sqlc.New(database)
	}
}
