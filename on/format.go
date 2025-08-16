package on

import (
	"github.com/razshare/frizzante/cli/path"
	"github.com/razshare/frizzante/cli/state"
	"github.com/razshare/frizzante/tui/messages"
	"os"
	"os/exec"
)

func Format() error {
	err := Touch()
	if err != nil {
		return err
	}

	gobin, err := path.Go(".")
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

	bunbin, err := path.Bun(*state.App)
	if err != nil {
		return err
	}

	pretty := exec.Command(bunbin, "x", "prettier", "--write", ".")
	pretty.Dir = *state.App
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
