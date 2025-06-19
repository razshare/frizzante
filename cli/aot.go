package cli

import (
	"embed"
	"github.com/razshare/frizzante/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type AotUtilities struct {
	efs embed.FS
}

func NewAotUtilities(efs embed.FS) *AotUtilities {
	return &AotUtilities{
		efs: efs,
	}
}

func (assets *AotUtilities) CreateOnDisk(to string) error {
	var from = "utilities"
	//var to = filepath.Join("frz/.generated", "utilities")
	var create func(from string, to string) error

	create = func(from string, to string) error {
		embeddedFiles, readDirError := assets.efs.ReadDir(from)
		if readDirError != nil {
			return readDirError
		}

		if !fs.IsDirectory(to) {
			mkdirError := os.MkdirAll(to, os.ModePerm)
			if mkdirError != nil {
				return mkdirError
			}
		}

		for _, embeddedFile := range embeddedFiles {
			fromFileName := strings.Join([]string{from, embeddedFile.Name()}, "/")
			toFileName := filepath.Join(to, strings.ReplaceAll(embeddedFile.Name(), "/", string(filepath.Separator)))
			if embeddedFile.IsDir() {
				createError := create(fromFileName, toFileName)
				if createError != nil {
					return createError
				}
				continue
			}

			embeddedContents, readError := assets.efs.ReadFile(fromFileName)
			if readError != nil {
				return readError
			}

			writeError := os.WriteFile(toFileName, embeddedContents, os.ModePerm)
			if writeError != nil {
				return writeError
			}
		}
		return nil
	}

	return create(from, to)
}

func CreateAotUtilitiesOnDisk(efs embed.FS, to string) {
	utilities := NewAotUtilities(efs)
	utilitiesError := utilities.CreateOnDisk(to)
	if utilitiesError != nil {
		log.Fatal(utilitiesError)
	}
}
