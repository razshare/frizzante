package action

import (
	"github.com/razshare/frizzante/tui/messages"
	"github.com/razshare/frizzante/tui/spinner"
	"os"
	"os/exec"
	"path/filepath"
)

func Install(o InstallOptions) error {
	err := Touch(TouchOptions{App: o.App})
	if err != nil {
		return err
	}

	s := spinner.New("installing go dependencies")

	go spinner.Start(s)
	tidy := exec.Command(o.Go, "mod", "tidy")
	tidy.Env = append(os.Environ())
	tidy.Stderr = os.Stderr
	tidy.Stdout = os.Stdout
	tidy.Stdin = os.Stdin
	err = tidy.Run()
	spinner.Stop(s)

	if err != nil {
		return err
	}

	bun, err := filepath.Rel(o.App, o.Bun)
	if err != nil {
		return err
	}

	ins := exec.Command(bun, "install")
	ins.Dir = o.App
	ins.Env = append(os.Environ())
	ins.Stderr = os.Stderr
	ins.Stdout = os.Stdout
	ins.Stdin = os.Stdin
	err = ins.Run()
	if err != nil {
		return err
	}

	messages.Success("project dependencies installed")

	return nil
}
