package action

import (
	"os"
	"os/exec"
	"path/filepath"

	"github.com/razshare/frizzante/files"
	"github.com/razshare/frizzante/internal/tui/messages"
	spinner2 "github.com/razshare/frizzante/internal/tui/spinner"
)

func Install(options InstallOptions) (err error) {
	if err = Touch(TouchOptions{App: options.App}); err != nil {
		return
	}

	spin := spinner2.New("installing go dependencies")

	go spinner2.Start(spin)
	tidy := exec.Command(options.Go, "mod", "tidy")
	tidy.Env = append(os.Environ())
	tidy.Stderr = os.Stderr
	tidy.Stdout = os.Stdout
	tidy.Stdin = os.Stdin
	err = tidy.Run()
	spinner2.Stop(spin)

	if err != nil {
		return
	}

	var bun string
	if files.IsFile(options.Bun) {
		if bun, err = filepath.Rel(options.App, options.Bun); err != nil {
			return err
		}
	} else if bun, err = exec.LookPath(options.Bun); err != nil {
		bun = options.Bun
	}

	ins := exec.Command(bun, "install")
	ins.Dir = options.App
	ins.Env = append(os.Environ())
	ins.Stderr = os.Stderr
	ins.Stdout = os.Stdout
	ins.Stdin = os.Stdin
	if err = ins.Run(); err != nil {
		return
	}

	messages.Success("project dependencies installed")

	return
}
