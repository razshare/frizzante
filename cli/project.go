package cli

import (
	"github.com/razshare/frizzante/files"
	"os"
	"path/filepath"
)

func OnCreateProject(prj string) {
	downloadError := files.DownloadFile("https://github.com/razshare/frizzante-starter/archive/refs/heads/main.zip", prj+".zip")
	if downloadError != nil {
		Fatal(downloadError)
	}

	unzipError := files.UnzipFile(prj+".zip", prj+".tmp")
	if unzipError != nil {
		Fatal(unzipError)
	}

	removeError := os.Remove(prj + ".zip")
	if removeError != nil {
		Fatal(removeError)
	}

	renameError := os.Rename(filepath.Join(prj+".tmp", "frizzante-starter-main"), prj)
	if renameError != nil {
		Fatal(renameError)
	}

	removeAllError := os.RemoveAll(filepath.Join(prj + ".tmp"))
	if removeAllError != nil {
		Fatal(removeAllError)
	}
}
