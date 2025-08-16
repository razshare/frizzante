package on

import (
	"embed"
	"github.com/razshare/frizzante/files"
	"github.com/razshare/frizzante/tui/messages"
	"os"
	"path/filepath"
)

func CreateProject(efs embed.FS, n string) error {
	err := files.DownloadFile("https://github.com/razshare/frizzante-starter/archive/refs/heads/main.zip", n+".zip")
	if err != nil {
		return err
	}

	err = files.UnzipFile(n+".zip", n+".tmp")
	if err != nil {
		return err
	}

	err = os.Remove(n + ".zip")
	if err != nil {
		return err
	}

	err = os.Rename(filepath.Join(n+".tmp", "frizzante-starter-main"), n)
	if err != nil {
		return err
	}

	err = os.RemoveAll(filepath.Join(n + ".tmp"))
	if err != nil {
		return err
	}

	messages.Successf("project created at %s", n)

	return Configure(efs, n)
}
