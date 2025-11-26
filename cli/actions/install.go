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

	spin := spinners.New("installing go packages")
	go spinners.Start(spin)
	if messages.Command(messages.CommandOptions{
		Environment: os.Environ(),
		Program:     options.Go,
		Args:        []string{"mod", "tidy"},
	}) {
		messages.Success("go packages installed")
	}
	spinners.Stop(spin)

	spin = spinners.New("installing javascript packages")
	go spinners.Start(spin)
	if messages.Command(messages.CommandOptions{
		DirectoryName: "app",
		Environment:   os.Environ(),
		Program:       bun,
		Args:          []string{"install"},
	}) {
		messages.Success("javascript packages installed")
	}
	spinners.Stop(spin)

	return
}
