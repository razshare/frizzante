package action

import (
	"os"
	"os/exec"
	"path/filepath"

	"github.com/razshare/frizzante/internal/project/lib/core/files"
	"github.com/razshare/frizzante/tui/messages"
	"github.com/razshare/frizzante/tui/spinner"
)

func Format(options FormatOptions) (err error) {
	spin := spinner.New("formatting code")
	go spinner.Start(spin)
	defer spinner.Stop(spin)

	if err = Touch(TouchOptions{App: options.App}); err != nil {
		return
	}

	gofmt := exec.Command(options.Go, "fmt", "./...")
	gofmt.Env = os.Environ()
	//gofmt.Stderr = os.Stderr
	//gofmt.Stdout = os.Stdout
	//gofmt.Stdin = os.Stdin
	if err = gofmt.Run(); err != nil {
		if gofmt.Err != nil {
			messages.Error(gofmt.Err.Error())
		}
		return
	}

	var bun string
	if files.IsFile(options.Bun) {
		if bun, err = filepath.Rel(options.App, options.Bun); err != nil {
			return
		}
	} else if bun, err = exec.LookPath(options.Bun); err != nil {
		bun = options.Bun
	}

	pretty := exec.Command(bun, "x", "prettier", "--write", ".")
	pretty.Dir = options.App
	pretty.Env = os.Environ()
	//pretty.Stderr = os.Stderr
	//pretty.Stdout = os.Stdout
	//pretty.Stdin = os.Stdin
	if err = pretty.Run(); err != nil {
		if pretty.Err != nil {
			messages.Error(pretty.Err.Error())
		}
		return
	}

	messages.Success("project formatted")

	return
}
