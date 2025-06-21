package cli

import (
	"embed"
	"github.com/razshare/frizzante/fs"
	"os"
	"path/filepath"
	"strings"
)

type Utilities struct {
	efs embed.FS
}

func NewUtilities(efs embed.FS) *Utilities {
	return &Utilities{
		efs: efs,
	}
}

func (utilities *Utilities) CreateOnDisk(to string) error {
	var from = "utilities"
	var create func(from string, to string) error

	create = func(from string, to string) error {
		embeddedFiles, readDirError := utilities.efs.ReadDir(from)
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

			embeddedContents, readError := utilities.efs.ReadFile(fromFileName)
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
