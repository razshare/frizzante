package action

import (
	"github.com/razshare/frizzante/cli/codegen"
	"github.com/razshare/frizzante/tui/messages"
	"os"
	"path/filepath"
)

func CreateProject(o CreateProjectOptions) error {
	url := "https://github.com/razshare/frizzante-starter/archive/refs/heads/main.zip"

	install, evict, err := codegen.Download(codegen.DownloadOptions{
		Url:  url,
		Auto: o.Auto,
	})

	if err != nil {
		return err
	}

	_, err = install(o.Project + ".tmp")
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

	err = evict()
	if err != nil {
		return err
	}

	messages.Successf("project has been created at %s", o.Project)
	messages.Tipf("change directory into %s and run frizzante --configure", o.Project)

	return nil
}
