package embeds

import (
	"embed"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/razshare/frizzante/v2/internal/project/lib/core/files"
)

func CopyFile(efs embed.FS, from string, to string) (err error) {
	var fromFile fs.File
	if fromFile, err = efs.Open(from); err != nil {
		return
	}
	defer func() {
		if cerr := fromFile.Close(); cerr != nil {
			if err == nil {
				err = cerr
			}
		}
	}()
	dir := filepath.Dir(to)
	if !files.IsDirectory(dir) {
		if err = os.MkdirAll(dir, os.ModePerm); err != nil {
			return
		}
	}
	if files.IsFile(to) {
		if err = os.Remove(to); err != nil {
			return
		}
	}
	var toFile *os.File
	if toFile, err = os.Create(to); err != nil {
		return
	}
	defer func() {
		if cerr := toFile.Close(); cerr != nil {
			if err == nil {
				err = cerr
			}
		}
	}()
	if _, err = io.Copy(toFile, fromFile); err != nil {
		return
	}
	return nil
}
