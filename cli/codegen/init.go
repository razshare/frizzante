package codegen

import (
	"embed"
	"github.com/razshare/frizzante/cli/user"
	"github.com/razshare/frizzante/embeds"
	"github.com/razshare/frizzante/files"
	"os"
	"path/filepath"
)

func Init(efs embed.FS) error {
	home, err := user.FrizzanteHome()
	if err != nil {
		return err
	}

	if !files.IsDirectory(filepath.Join(home, "template", "lib")) {
		err = embeds.CopyDirectory(efs, "template/lib", filepath.Join(home, "template", "lib"))
		if err != nil {
			return err
		}
	}

	if !files.IsDirectory(filepath.Join(home, "template", "project")) {
		if files.IsFile(filepath.Join(home, "template", "project.zip")) {
			err = os.RemoveAll(filepath.Join(home, "template", "project.zip"))
			if err != nil {
				return err
			}
		}
		err = embeds.CopyFile(efs, "template/project.zip", filepath.Join(home, "template", "project.zip"))
		if err != nil {
			return err
		}

		err = files.UnzipFile(filepath.Join(home, "template", "project.zip"), filepath.Join(home, "template", "project.base"))
		if err != nil {
			return err
		}

		err = os.Rename(filepath.Join(home, "template", "project.base", "project.tmp"), filepath.Join(home, "template", "project"))
		if err != nil {
			return err
		}

		err = os.RemoveAll(filepath.Join(home, "template", "project.base"))
		if err != nil {
			return err
		}

		err = os.RemoveAll(filepath.Join(home, "template", "project.zip"))
		if err != nil {
			return err
		}
	}

	return nil
}
