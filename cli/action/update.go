package action

import (
	"github.com/razshare/frizzante/tui/messages"
	"github.com/razshare/frizzante/tui/spinner"
	"os"
	"os/exec"
	"path/filepath"
)

func Update(opts UpdateOptions) error {
	err := Touch(TouchOptions{App: opts.App})
	if err != nil {
		return err
	}

	spin := spinner.New("updating go dependencies")

	go spinner.Start(spin)
	get := exec.Command(opts.Go, "get", "-u", "./...")
	get.Env = append(os.Environ())
	get.Stderr = os.Stderr
	get.Stdout = os.Stdout
	get.Stdin = os.Stdin
	err = get.Run()
	spinner.Stop(spin)

	if err != nil {
		return err
	}

	bun, err := filepath.Rel(opts.App, opts.Bun)
	if err != nil {
		return err
	}

	pretty := exec.Command(bun, "update")
	pretty.Dir = opts.App
	pretty.Env = append(os.Environ())
	pretty.Stderr = os.Stderr
	pretty.Stdout = os.Stdout
	pretty.Stdin = os.Stdin
	err = pretty.Run()
	if err != nil {
		return err
	}

	messages.Success("project dependencies updated")

	return nil
}
