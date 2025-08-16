package on

import (
	"github.com/razshare/frizzante/cli"
	"github.com/razshare/frizzante/cli/path"
	"github.com/razshare/frizzante/tui/messages"
	"github.com/razshare/frizzante/tui/spinner"
	"os"
	"os/exec"
)

func Update(c *cli.Cli) error {
	err := Touch(c)
	if err != nil {
		return err
	}

	s := spinner.New("updating go dependencies")

	gobin, err := path.Go(c, ".")
	if err != nil {
		return err
	}

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

	bunbin, err := path.Bun(c, *c.Flags.App)
	if err != nil {
		return err
	}

	pretty := exec.Command(bunbin, "update")
	pretty.Dir = *c.Flags.App
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
