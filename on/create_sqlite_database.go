package on

import (
	"embed"
	"github.com/razshare/frizzante/files"
	"github.com/razshare/frizzante/tui/confirm"
	"github.com/razshare/frizzante/tui/messages"
	"io"
	"os"
)

func CreateSqliteDatabase(efs embed.FS) {
	if files.IsFile("database.sqlite") {
		if confirm.Send(true, "database.sqlite already exists, would you like to overwrite it?") {
			err := os.Remove("database.sqlite")
			if err != nil {
				messages.Fatal(err)
			}
		}
	}

	efile, err := efs.Open("database.sqlite")
	if err != nil {
		messages.Fatal(err)
	}

	file, err := os.Create("database.sqlite")
	if err != nil {
		messages.Fatal(err)
	}

	_, err = io.Copy(file, efile)
	if err != nil {
		messages.Fatal(err)
	}

	messages.Success("database.sqlite created")
}
