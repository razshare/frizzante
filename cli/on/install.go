package on

import (
	"github.com/razshare/frizzante/tui/messages"
	"github.com/razshare/frizzante/tui/spinner"
	"os"
	"os/exec"
)

func Install(app string, gobin string, bunbin string) error {
	err := Touch(app)
	if err != nil {
		return err
	}

	s := spinner.New("installing go dependencies")

	go spinner.Start(s)
	tidy := exec.Command(gobin, "mod", "tidy")
	tidy.Env = append(os.Environ())
	tidy.Stderr = os.Stderr
	tidy.Stdout = os.Stdout
	tidy.Stdin = os.Stdin
	err = tidy.Run()
	spinner.Stop(s)

	if err != nil {
		return err
	}

	ins := exec.Command(bunbin, "install")
	ins.Dir = app
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
