package action

import (
	"github.com/razshare/frizzante/files"
	"os"
	"path/filepath"
)

func Touch(o TouchOptions) error {
	touch := func(n string) error {
		dn := filepath.Dir(n)

		if !files.IsDirectory(dn) {
			err := os.MkdirAll(dn, os.ModePerm)
			if err != nil {
				return err
			}
		}

		file, err := os.OpenFile(n, os.O_RDONLY|os.O_CREATE, 0666)
		if err != nil {
			return err
		}

		err = file.Close()
		if err != nil {
			return err
		}
		return nil
	}

	err := os.MkdirAll(filepath.Join(o.App, "dist"), os.ModePerm)
	if err != nil {
		return err
	}

	err = touch(filepath.Join(o.App, "dist", "server.js"))
	if err != nil {
		return err
	}

	return touch(filepath.Join(o.App, "dist", "client", "index.html"))
}
