package on

import (
	"github.com/razshare/frizzante/files"
	"github.com/razshare/frizzante/tui/messages"
	"os"
	"path/filepath"
)

func CreateProject(prj string) {
	downloadError := files.DownloadFile("https://github.com/razshare/frizzante-starter/archive/refs/heads/main.zip", prj+".zip")
	if downloadError != nil {
		messages.Fatal(downloadError)
	}

	unzipError := files.UnzipFile(prj+".zip", prj+".tmp")
	if unzipError != nil {
		messages.Fatal(unzipError)
	}

	removeError := os.Remove(prj + ".zip")
	if removeError != nil {
		messages.Fatal(removeError)
	}

	renameError := os.Rename(filepath.Join(prj+".tmp", "frizzante-starter-main"), prj)
	if renameError != nil {
		messages.Fatal(renameError)
	}

	removeAllError := os.RemoveAll(filepath.Join(prj + ".tmp"))
	if removeAllError != nil {
		messages.Fatal(removeAllError)
	}
}
