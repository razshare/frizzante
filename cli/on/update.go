package on

import (
	"github.com/razshare/frizzante/tui/messages"
	"github.com/razshare/frizzante/tui/spinner"
	"os"
	"os/exec"
)

func Update(app string, gobin string, bunbin string) error {
	err := Touch(app)
	if err != nil {
		return err
	}

	s := spinner.New("updating go dependencies")

	go spinner.Start(s)
	get := exec.Command(gobin, "get", "-u", "./...")
	get.Env = append(os.Environ())
	get.Stderr = os.Stderr
	get.Stdout = os.Stdout
	get.Stdin = os.Stdin
	err = get.Run()
	spinner.Stop(s)

	if err != nil {
		return err
	}

	pretty := exec.Command(bunbin, "update")
	pretty.Dir = app
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
