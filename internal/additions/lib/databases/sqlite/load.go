package databases

import (
	"database/sql"
	"os"

	_ "github.com/mattn/go-sqlite3"
	"github.com/razshare/frizzante/internal/project/lib/core/files"
)

func Load() (database *sql.DB, err error) {
	if !files.IsFile("source.sqlite") {
		var data []byte
		if data, err = Efs.ReadFile("source.sqlite"); err != nil {
			return
		}

		if err = os.WriteFile("source.sqlite", data, os.ModePerm); err != nil {
			return
		}
	}

	if database, err = sql.Open("sqlite3", "file:source.sqlite?cache=shared"); err != nil {
		return
	}

	return
}
