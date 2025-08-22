package action

import (
	"github.com/razshare/frizzante/files"
	"os"
	"path/filepath"
)

func Touch(opts TouchOptions) error {
	touch := func(name string) error {
		dir := filepath.Dir(name)

		if !files.IsDirectory(dir) {
			err := os.MkdirAll(dir, os.ModePerm)
			if err != nil {
				return err
			}
		}

		file, err := os.OpenFile(name, os.O_RDONLY|os.O_CREATE, 0666)
		if err != nil {
			return err
		}

		err = file.Close()
		if err != nil {
			return err
		}
		return nil
	}

	err := os.MkdirAll(filepath.Join(opts.App, "dist"), os.ModePerm)
	if err != nil {
		return err
	}

	err = touch(filepath.Join(opts.App, "dist", "server.js"))
	if err != nil {
		return err
	}

	return touch(filepath.Join(opts.App, "dist", "client", "index.html"))
}
