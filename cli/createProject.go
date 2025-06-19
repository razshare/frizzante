package cli

import (
	"embed"
	"github.com/pterm/pterm"
	"github.com/razshare/frizzante/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func CreateProject(efs embed.FS) {
	if !*FlagCreateProject {
		return
	}

	*FlagOut = "asd"

	var showError error

	if "" == *FlagOut {
		*FlagOut, showError = pterm.
			DefaultInteractiveTextInput.
			Show("The name of the project is")
		if showError != nil {
			log.Fatal(showError)
		}
	}

	var list func(efs embed.FS, directoryName string, callback func(efs embed.FS, fileName string, data []byte))
	list = func(efs embed.FS, directoryName string, callback func(efs embed.FS, fileName string, data []byte)) {
		if strings.HasSuffix(directoryName, "/") {
			directoryName = strings.TrimSuffix(directoryName, "/")
		} else if strings.HasSuffix(directoryName, "\\") {
			directoryName = strings.TrimSuffix(directoryName, "\\")
		}

		entries, readDirError := efs.ReadDir(directoryName)
		if readDirError != nil {
			log.Fatal(readDirError)
		}

		for _, entry := range entries {
			if entry.IsDir() {
				list(efs, directoryName+"/"+entry.Name(), callback)
				continue
			}
			fileName := filepath.Join(
				directoryName,
				strings.ReplaceAll(entry.Name(), "/", string(filepath.Separator)),
			)
			efsFileName := directoryName + "/" + entry.Name()
			data, readError := efs.ReadFile(efsFileName)
			if readError != nil {
				log.Fatal(readError)
			}
			callback(efs, fileName, data)
		}
	}

	list(efs, "frizzante-starter", func(efs embed.FS, fileName string, data []byte) {
		directoryName := filepath.Dir(fileName)

		if !fs.IsDirectory(directoryName) {
			mkdirError := os.MkdirAll(directoryName, os.ModePerm)
			if mkdirError != nil {
				log.Fatal(mkdirError)
			}
		}

		writeError := os.WriteFile(fileName, data, os.ModePerm)
		if writeError != nil {
			log.Fatal(writeError)
		}
	})
}
