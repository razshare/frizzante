package cli

import (
	"github.com/razshare/frizzante/files"
	"os"
	"path/filepath"
)

func OnTouch() {
	touch := func(fileName string) {
		directoryName := filepath.Dir(fileName)

		if !files.IsDirectory(directoryName) {
			mkdirAllError := os.MkdirAll(directoryName, os.ModePerm)
			if mkdirAllError != nil {
				Fatal(mkdirAllError)
			}
		}

		file, openError := os.OpenFile(fileName, os.O_RDONLY|os.O_CREATE, 0666)
		if openError != nil {
			Fatal(openError)
		}

		closeError := file.Close()
		if closeError != nil {
			Fatal(closeError)
		}
	}

	mkdirError := os.MkdirAll(filepath.Join("app", "dist"), os.ModePerm)
	if mkdirError != nil {
		Fatal(mkdirError)
	}

	touch(filepath.Join("app", "dist", "server.js"))
	touch(filepath.Join("app", "dist", "client", "index.html"))
}
