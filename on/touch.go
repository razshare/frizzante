package on

import (
	"github.com/razshare/frizzante/cli/flags"
	"github.com/razshare/frizzante/files"
	"github.com/razshare/frizzante/tui/messages"
	"os"
	"path/filepath"
)

func Touch() {
	touch := func(fileName string) {
		directoryName := filepath.Dir(fileName)

		if !files.IsDirectory(directoryName) {
			mkdirAllError := os.MkdirAll(directoryName, os.ModePerm)
			if mkdirAllError != nil {
				messages.Fatal(mkdirAllError)
			}
		}

		file, openError := os.OpenFile(fileName, os.O_RDONLY|os.O_CREATE, 0666)
		if openError != nil {
			messages.Fatal(openError)
		}

		closeError := file.Close()
		if closeError != nil {
			messages.Fatal(closeError)
		}
	}

	mkdirError := os.MkdirAll(filepath.Join(*flags.App, "dist"), os.ModePerm)
	if mkdirError != nil {
		messages.Fatal(mkdirError)
	}

	touch(filepath.Join(*flags.App, "dist", "server.js"))
	touch(filepath.Join(*flags.App, "dist", "client", "index.html"))
}
