package sqlite

import (
	"database/sql"
	"embed"
	_ "github.com/mattn/go-sqlite3"
	"github.com/razshare/frizzante/files"
	"github.com/razshare/frizzante/template/lib/database/sqlite/sqlc"
	"log"
	"os"
	//gen:mod "github.com/razshare/frizzante/template" "main"
)

//gen:mod "queries" "Queries"
var queries *sqlc.Queries

//go:embed source.sqlite
//gen:mod "efs" "Efs"
var dbf embed.FS

func init() {
	if !files.IsFile("source.sqlite") {
		//gen:mod "efs" "Efs"
		data, readError := dbf.ReadFile("source.sqlite")
		if readError != nil {
			log.Fatal(readError)
		}
		writeError := os.WriteFile("source.sqlite", data, os.ModePerm)
		if writeError != nil {
			log.Fatal(writeError)
		}
	}

	db, err := sql.Open("sqlite3", "file:source.sqlite?cache=shared")
	if err != nil {
		log.Fatal(err)
	}

	//gen:mod "queries" "Queries"
	queries = sqlc.New(db)
}
