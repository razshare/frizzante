package database

import (
	"database/sql"
	"embed"
	_ "github.com/mattn/go-sqlite3"
	"github.com/razshare/frizzante/files"
	//gen:mod "github.com/razshare/frizzante/template" "main"
	"github.com/razshare/frizzante/template/lib/database/sqlc"
	"log"
	"os"
)

//go:embed source.sqlite
var efs embed.FS

//gen:mod "queries" "Queries"
var queries *sqlc.Queries

func init() {
	if !files.IsFile("source.sqlite") {
		d, err := efs.ReadFile("source.sqlite")
		if err != nil {
			log.Fatal(err)
		}
		err = os.WriteFile("source.sqlite", d, os.ModePerm)
		if err != nil {
			log.Fatal(err)
		}
	}

	db, err := sql.Open("sqlite3", "file:source.sqlite?cache=shared")
	if err != nil {
		log.Fatal(err)
	}

	//gen:mod "queries" "Queries"
	queries = sqlc.New(db)
}
