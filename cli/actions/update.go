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
		return
	}

	spin := spinners.New("updating go packages")
	go spinners.Start(spin)
	if messages.Command(messages.CommandOptions{
		Env:  os.Environ(),
		Name: options.Go,
		Args: []string{"get", "-u", "./..."},
	}) {
		spinners.Stop(spin)
		messages.Success("go packages updated")
	}

	spin = spinners.New("updating javascript packages")
	go spinners.Start(spin)
	if messages.Command(messages.CommandOptions{
		Dir:  "app",
		Env:  os.Environ(),
		Name: bun,
		Args: []string{"update"},
	}) {
		spinners.Stop(spin)
		messages.Success("javascript packages updated")
	}

	return
}
