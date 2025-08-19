package action

import (
	"github.com/razshare/frizzante/files"
	"github.com/razshare/frizzante/tui/messages"
	"os"
	"path/filepath"
)

func CreateProject(o CreateProjectOptions) error {
	err := files.DownloadFile("https://github.com/razshare/frizzante-starter/archive/refs/heads/main.zip", o.Project+".zip")
	if err != nil {
		return err
	}

	err = files.UnzipFile(o.Project+".zip", o.Project+".tmp")
	if err != nil {
		return err
	}

	err = os.Remove(o.Project + ".zip")
	if err != nil {
		return err
	}

	err = os.Rename(filepath.Join(o.Project+".tmp", "frizzante-starter-main"), o.Project)
	if err != nil {
		return err
	}

	err = os.RemoveAll(filepath.Join(o.Project + ".tmp"))
	if err != nil {
		return err
	}

	messages.Successf("project created at %s", o.Project)

	return nil
}
