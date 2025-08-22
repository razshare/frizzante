package action

import (
	"github.com/razshare/frizzante/tui/messages"
	"github.com/razshare/frizzante/tui/spinner"
	"os"
	"os/exec"
	"path/filepath"
)

func Install(opts InstallOptions) error {
	err := Touch(TouchOptions{App: opts.App})
	if err != nil {
		return err
	}

	spin := spinner.New("installing go dependencies")

	go spinner.Start(spin)
	tidy := exec.Command(opts.Go, "mod", "tidy")
	tidy.Env = append(os.Environ())
	tidy.Stderr = os.Stderr
	tidy.Stdout = os.Stdout
	tidy.Stdin = os.Stdin
	err = tidy.Run()
	spinner.Stop(spin)

	if err != nil {
		return err
	}

	bun, err := filepath.Rel(opts.App, opts.Bun)
	if err != nil {
		return err
	}

	ins := exec.Command(bun, "install")
	ins.Dir = opts.App
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
