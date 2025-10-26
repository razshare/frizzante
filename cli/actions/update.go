package actions

import (
	"os"
	"os/exec"
	"path/filepath"

	"github.com/razshare/frizzante/internal/project/lib/core/files"
	"github.com/razshare/frizzante/tui/messages"
	"github.com/razshare/frizzante/tui/spinners"
)

func Update(options UpdateOptions) (err error) {
	spin := spinners.New("updating packages")
	go spinners.Start(spin)
	defer spinners.Stop(spin)

	if err = Touch(TouchOptions{}); err != nil {
		return
	}
	var bun string
	if files.IsFile(options.Bun) {
		if bun, err = filepath.Rel("app", options.Bun); err != nil {
			return
		}
	} else if bun, err = exec.LookPath(options.Bun); err != nil {
		bun = options.Bun
	}

	if messages.Command("", os.Environ(), options.Go, "get", "-u", "./...") {
		messages.Success("go packages updated")
	}

	if messages.Command("app", os.Environ(), bun, "update") {
		messages.Success("js packages updated")
	}

	return
}
