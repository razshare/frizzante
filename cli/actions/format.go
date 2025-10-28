package actions

import (
	"os"
	"os/exec"
	"path/filepath"

	"github.com/razshare/frizzante/internal/project/lib/core/files"
	"github.com/razshare/frizzante/tui/messages"
	"github.com/razshare/frizzante/tui/spinners"
)

func Format(options FormatOptions) (err error) {
	spin := spinners.New("formatting code")
	go spinners.Start(spin)
	defer spinners.Stop(spin)

	if err = Touch(TouchOptions{}); err != nil {
		return
	}

	gofmt := exec.Command(options.Go, "fmt", "./...")
	gofmt.Env = os.Environ()
	gofmt.Stderr = os.Stderr
	gofmt.Stdout = os.Stdout
	gofmt.Stdin = os.Stdin
	if err = gofmt.Run(); err != nil {
		if gofmt.Err != nil {
			messages.Error(gofmt.Err.Error())
		}
		return
	}

	var bun string
	if files.IsFile(options.Bun) {
		if bun, err = filepath.Rel("app", options.Bun); err != nil {
			return
		}
	} else if bun, err = exec.LookPath(options.Bun); err != nil {
		bun = options.Bun
		return
	}

	pretty := exec.Command(bun, "x", "prettier", "--write", ".")
	pretty.Dir = "app"
	pretty.Env = os.Environ()
	pretty.Stderr = os.Stderr
	pretty.Stdout = os.Stdout
	pretty.Stdin = os.Stdin
	if err = pretty.Run(); err != nil {
		spinners.Stop(spin)
		if pretty.Err != nil {
			messages.Error(pretty.Err.Error())
		}
		return
	}

	messages.Success("project formatted")

	return
}
