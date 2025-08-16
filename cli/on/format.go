package on

import (
	"github.com/razshare/frizzante/cli"
	"github.com/razshare/frizzante/cli/path"
	"github.com/razshare/frizzante/tui/messages"
	"os"
	"os/exec"
)

func Format(c *cli.Cli) error {
	err := Touch(c)
	if err != nil {
		return err
	}

	gobin, err := path.Go(c, ".")
	if err != nil {
		return err
	}

	gofmt := exec.Command(gobin, "fmt")
	gofmt.Env = append(os.Environ())
	gofmt.Stderr = os.Stderr
	gofmt.Stdout = os.Stdout
	gofmt.Stdin = os.Stdin
	err = gofmt.Run()
	if err != nil {
		return err
	}

	bunbin, err := path.Bun(c, *c.Flags.App)
	if err != nil {
		return err
	}

	pretty := exec.Command(bunbin, "x", "prettier", "--write", ".")
	pretty.Dir = *c.Flags.App
	pretty.Env = append(os.Environ())
	pretty.Stderr = os.Stderr
	pretty.Stdout = os.Stdout
	pretty.Stdin = os.Stdin
	err = pretty.Run()
	if err != nil {
		return err
	}

	messages.Success("project formatted")

	return nil
}
