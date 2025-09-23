package action

import (
	"os"
	"os/exec"
	"path/filepath"

	"github.com/razshare/frizzante/internal/project/lib/core/files"
	"github.com/razshare/frizzante/tui/messages"
	"github.com/razshare/frizzante/tui/spinner"
)

func Update(options UpdateOptions) (err error) {
	spin := spinner.New("updating packages")
	go spinner.Start(spin)
	defer spinner.Stop(spin)

	if err = Touch(TouchOptions{App: options.App}); err != nil {
		return
	}

	get := exec.Command(options.Go, "get", "-u", "./...")
	get.Env = os.Environ()
	//get.Stderr = os.Stderr
	//get.Stdout = os.Stdout
	//get.Stdin = os.Stdin
	err = get.Run()

	if err != nil {
		if get.Err != nil {
			messages.Error(get.Err.Error())
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

	pretty := exec.Command(bun, "update")
	pretty.Dir = options.App
	pretty.Env = os.Environ()
	//pretty.Stderr = os.Stderr
	//pretty.Stdout = os.Stdout
	//pretty.Stdin = os.Stdin
	err = pretty.Run()

	if err != nil {
		if pretty.Err != nil {
			messages.Error(pretty.Err.Error())
		}
		return
	}

	messages.Success("project dependencies updated")

	return
}
