package cli

import (
	"embed"
	"github.com/razshare/frizzante/files"
	"io"
	"os"
)

func OnSqliteDatabase(efs embed.FS) {
	if files.IsFile("database.sqlite") {
		if Confirm("database.sqlite already exists, would you like to overwrite it?") {
			err := os.Remove("database.sqlite")
			if err != nil {
				Fatal(err)
			}
		}
	}

	efile, err := efs.Open("database.sqlite")
	if err != nil {
		Fatal(err)
	}

	file, err := os.Create("database.sqlite")
	if err != nil {
		Fatal(err)
	}

	_, err = io.Copy(file, efile)
	if err != nil {
		Fatal(err)
	}

	Success("database.sqlite created")
}
