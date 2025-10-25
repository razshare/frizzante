package actions

import (
	"os"
	"os/exec"
	"path/filepath"

	"github.com/razshare/frizzante/internal/project/lib/core/files"
	"github.com/razshare/frizzante/tui/messages"
	"github.com/razshare/frizzante/tui/spinners"
)

func Install(options InstallOptions) (err error) {
	spin := spinners.New("installing packages")
	go spinners.Start(spin)
	defer spinners.Stop(spin)

	if err = Touch(TouchOptions{}); err != nil {
		return
	}

	var bun string
	if files.IsFile(options.Bun) {
		if bun, err = filepath.Rel("app", options.Bun); err != nil {
			return err
		}
	} else if bun, err = exec.LookPath(options.Bun); err != nil {
		bun = options.Bun
	}

	if messages.Command(".", os.Environ(), options.Go, "mod", "tidy") {
		messages.Success("go packages installed")
	}

	if messages.Command("app", os.Environ(), bun, "install") {
		messages.Success("js packages installed")
	}

	return
}
