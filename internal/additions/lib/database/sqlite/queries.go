package database

import (
	"database/sql"
	"embed"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
	"github.com/razshare/frizzante/internal/additions/lib/database/sqlite/sqlc"
	"github.com/razshare/frizzante/internal/project/lib/core/files"
)

var Queries *sqlc.Queries

//go:embed source.sqlite
var Efs embed.FS

func init() {
	if !files.IsFile("source.sqlite") {
		data, err := Efs.ReadFile("source.sqlite")
		if err != nil {
			log.Fatal(err)
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

	Queries = sqlc.New(db)
}
