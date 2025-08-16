package on

import (
	"github.com/razshare/frizzante/cli"
	"github.com/razshare/frizzante/files"
	"github.com/razshare/frizzante/tui/messages"
	"os"
	"path/filepath"
)

func CreateProject(c *cli.Cli, dst string) error {
	err := files.DownloadFile("https://github.com/razshare/frizzante-starter/archive/refs/heads/main.zip", dst+".zip")
	if err != nil {
		return err
	}

	err = files.UnzipFile(dst+".zip", dst+".tmp")
	if err != nil {
		return err
	}

	err = os.Remove(dst + ".zip")
	if err != nil {
		return err
	}

	err = os.Rename(filepath.Join(dst+".tmp", "frizzante-starter-main"), dst)
	if err != nil {
		return err
	}

	err = os.RemoveAll(filepath.Join(dst + ".tmp"))
	if err != nil {
		return err
	}

	messages.Successf("project created at %s", dst)

	return Configure(c, dst)
}
