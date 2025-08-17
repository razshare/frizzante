package action

import (
	"github.com/razshare/frizzante/tui/messages"
	"os"
	"os/exec"
	"path/filepath"
)

func Format(o FormatOptions) error {
	err := Touch(TouchOptions{App: o.App})
	if err != nil {
		return err
	}

	gofmt := exec.Command(o.Go, "fmt", "./...")
	gofmt.Env = append(os.Environ())
	gofmt.Stderr = os.Stderr
	gofmt.Stdout = os.Stdout
	gofmt.Stdin = os.Stdin
	err = gofmt.Run()
	if err != nil {
		return err
	}

	bun, err := filepath.Rel(o.App, o.Bun)
	if err != nil {
		return err
	}

	pretty := exec.Command(bun, "x", "prettier", "--write", ".")
	pretty.Dir = o.App
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
