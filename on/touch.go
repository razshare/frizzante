package on

import (
	"github.com/razshare/frizzante/cli/flags"
	"github.com/razshare/frizzante/files"
	"github.com/razshare/frizzante/tui/messages"
	"os"
	"path/filepath"
)

func Touch() {
	touch := func(n string) {
		dn := filepath.Dir(n)

		if !files.IsDirectory(dn) {
			err := os.MkdirAll(dn, os.ModePerm)
			if err != nil {
				messages.Fatal(err)
			}
		}

		file, err := os.OpenFile(n, os.O_RDONLY|os.O_CREATE, 0666)
		if err != nil {
			messages.Fatal(err)
		}

		err = file.Close()
		if err != nil {
			messages.Fatal(err)
		}
	}

	err := os.MkdirAll(filepath.Join(*flags.App, "dist"), os.ModePerm)
	if err != nil {
		messages.Fatal(err)
	}

	touch(filepath.Join(*flags.App, "dist", "server.js"))
	touch(filepath.Join(*flags.App, "dist", "client", "index.html"))
}
