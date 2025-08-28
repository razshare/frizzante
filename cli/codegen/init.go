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

	if files.IsDirectory(filepath.Join(home, "internal", "template")) {
		if err = os.RemoveAll(filepath.Join(home, "internal", "template")); err != nil {
			return
		}
	}

	if err = embeds.CopyDirectory(efs, "internal/template/lib", filepath.Join(home, "internal", "template", "lib")); err != nil {
		return
	}

	if files.IsFile(filepath.Join(home, "internal", "template", "project.zip")) {
		if err = os.RemoveAll(filepath.Join(home, "internal", "template", "project.zip")); err != nil {
			return
		}
	}

	if err = embeds.CopyFile(efs, "internal/template/project.zip", filepath.Join(home, "internal", "template", "project.zip")); err != nil {
		return
	}

	if err = files.UnzipFile(filepath.Join(home, "internal", "template", "project.zip"), filepath.Join(home, "internal", "template", "project.base")); err != nil {
		return
	}

	if err = os.Rename(filepath.Join(home, "internal", "template", "project.base", "project.tmp"), filepath.Join(home, "internal", "template", "project")); err != nil {
		return
	}

	if err = os.RemoveAll(filepath.Join(home, "internal", "template", "project.base")); err != nil {
		return
	}

	if err = os.RemoveAll(filepath.Join(home, "internal", "template", "project.zip")); err != nil {
		return
	}

	return
}
