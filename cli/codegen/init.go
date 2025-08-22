package codegen

import (
	"embed"
	"github.com/razshare/frizzante/cli/user"
	"github.com/razshare/frizzante/embeds"
	"github.com/razshare/frizzante/files"
	"os"
	"path/filepath"
)

func Init(efs embed.FS) (err error) {
	var home string
	if home, err = user.FrizzanteHome(); err != nil {
		return
	}

	if files.IsDirectory(filepath.Join(home, "template")) {
		if err = os.RemoveAll(filepath.Join(home, "template")); err != nil {
			return
		}
	}

	if err = embeds.CopyDirectory(efs, "template/lib", filepath.Join(home, "template", "lib")); err != nil {
		return
	}

	if files.IsFile(filepath.Join(home, "template", "project.zip")) {
		if err = os.RemoveAll(filepath.Join(home, "template", "project.zip")); err != nil {
			return
		}
	}

	if err = embeds.CopyFile(efs, "template/project.zip", filepath.Join(home, "template", "project.zip")); err != nil {
		return
	}

	if err = files.UnzipFile(filepath.Join(home, "template", "project.zip"), filepath.Join(home, "template", "project.base")); err != nil {
		return
	}

	if err = os.Rename(filepath.Join(home, "template", "project.base", "project.tmp"), filepath.Join(home, "template", "project")); err != nil {
		return
	}

	if err = os.RemoveAll(filepath.Join(home, "template", "project.base")); err != nil {
		return
	}

	if err = os.RemoveAll(filepath.Join(home, "template", "project.zip")); err != nil {
		return
	}

	return
}
