package on

import (
	"github.com/razshare/frizzante/tui/messages"
	"os"
	"os/exec"
)

func Format(app string, gobin string, bunbin string) error {
	err := Touch(app)
	if err != nil {
		return err
	}

	gofmt := exec.Command(gobin, "fmt", "./...")
	gofmt.Env = append(os.Environ())
	gofmt.Stderr = os.Stderr
	gofmt.Stdout = os.Stdout
	gofmt.Stdin = os.Stdin
	err = gofmt.Run()
	if err != nil {
		return err
	}

	pretty := exec.Command(bunbin, "x", "prettier", "--write", ".")
	pretty.Dir = app
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
