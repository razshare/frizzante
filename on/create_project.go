package on

import (
	"github.com/razshare/frizzante/files"
	"github.com/razshare/frizzante/tui/messages"
	"os"
	"path/filepath"
)

func CreateProject(prj string) {
	err := files.DownloadFile("https://github.com/razshare/frizzante-starter/archive/refs/heads/main.zip", prj+".zip")
	if err != nil {
		messages.Fatal(err)
	}

	err = files.UnzipFile(prj+".zip", prj+".tmp")
	if err != nil {
		messages.Fatal(err)
	}

	err = os.Remove(prj + ".zip")
	if err != nil {
		messages.Fatal(err)
	}

	err = os.Rename(filepath.Join(prj+".tmp", "frizzante-starter-main"), prj)
	if err != nil {
		messages.Fatal(err)
	}

	err = os.RemoveAll(filepath.Join(prj + ".tmp"))
	if err != nil {
		messages.Fatal(err)
	}
}
