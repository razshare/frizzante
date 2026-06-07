package databases

import (
	"database/sql"
	"log"

	"github.com/razshare/frizzante/internal/project/lib/core/databases/sqlc"
)

func init() {
	var err error
	var database *sql.DB
	if database, err = Load(); err != nil {
		log.Fatal(err)
	}
	var queriesLocal *sqlc.Queries
	if queriesLocal = sqlc.New(database); queriesLocal == nil {
		log.Fatal("could not construct database queries object")
	}
	Queries = *queriesLocal
}
