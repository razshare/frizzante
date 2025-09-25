package action

import (
	"os"
	"os/exec"
	"path/filepath"

	"github.com/razshare/frizzante/internal/project/lib/core/files"
	"github.com/razshare/frizzante/tui/messages"
	"github.com/razshare/frizzante/tui/spinner"
)

func Update(options UpdateOptions) (err error) {
	spin := spinner.New("updating packages")
	go spinner.Start(spin)
	defer spinner.Stop(spin)

	if err = Touch(TouchOptions{App: options.App}); err != nil {
		return
	}
	var bun string
	if files.IsFile(options.Bun) {
		if bun, err = filepath.Rel(options.App, options.Bun); err != nil {
			return
		}
	} else if bun, err = exec.LookPath(options.Bun); err != nil {
		bun = options.Bun
	}

	if messages.Command(".", os.Environ(), options.Go, "get", "-u", "./...") {
		messages.Success("go packages updated")
	}

	if messages.Command(options.App, os.Environ(), bun, "update") {
		messages.Success("js packages updated")
	}

	return
}
