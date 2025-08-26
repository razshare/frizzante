package action

import (
	"github.com/razshare/frizzante/files"
	"github.com/razshare/frizzante/tui/messages"
	"os"
	"os/exec"
	"path/filepath"
)

func Format(opts FormatOptions) error {

	if err := Touch(TouchOptions{App: opts.App}); err != nil {
		return err
	}

	gofmt := exec.Command(opts.Go, "fmt", "./...")
	gofmt.Env = append(os.Environ())
	gofmt.Stderr = os.Stderr
	gofmt.Stdout = os.Stdout
	gofmt.Stdin = os.Stdin
	if err := gofmt.Run(); err != nil {
		return err
	}

	var bun string
	var err error
	if files.IsFile(opts.Bun) {
		if bun, err = filepath.Rel(opts.App, opts.Bun); err != nil {
			return err
		}
	} else if bun, err = exec.LookPath(opts.Bun); err != nil {
		bun = opts.Bun
	}

	pretty := exec.Command(bun, "x", "prettier", "--write", ".")
	pretty.Dir = opts.App
	pretty.Env = append(os.Environ())
	pretty.Stderr = os.Stderr
	pretty.Stdout = os.Stdout
	pretty.Stdin = os.Stdin
	if err := pretty.Run(); err != nil {
		return err
	}

	messages.Success("project formatted")

	return nil
}
