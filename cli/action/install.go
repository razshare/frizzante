package action

import (
	"os"
	"os/exec"
	"path/filepath"

	"github.com/razshare/frizzante/internal/project/lib/core/files"
	"github.com/razshare/frizzante/tui/messages"
	"github.com/razshare/frizzante/tui/spinner"
)

func Install(options InstallOptions) (err error) {
	if err = Touch(TouchOptions{App: options.App}); err != nil {
		return
	}

	spin := spinner.New("installing go dependencies")

	go spinner.Start(spin)
	tidy := exec.Command(options.Go, "mod", "tidy")
	tidy.Env = append(os.Environ())
	//tidy.Stderr = os.Stderr
	//tidy.Stdout = os.Stdout
	//tidy.Stdin = os.Stdin
	err = tidy.Run()
	spinner.Stop(spin)

	if err != nil {
		if tidy.Err != nil {
			messages.Error(tidy.Err.Error())
		}
		return
	}

	var bun string
	if files.IsFile(options.Bun) {
		if bun, err = filepath.Rel(options.App, options.Bun); err != nil {
			return err
		}
	} else if bun, err = exec.LookPath(options.Bun); err != nil {
		bun = options.Bun
	}

	install := exec.Command(bun, "install")
	install.Dir = options.App
	install.Env = append(os.Environ())
	//ins.Stderr = os.Stderr
	//ins.Stdout = os.Stdout
	//ins.Stdin = os.Stdin
	if err = install.Run(); err != nil {
		if install.Err != nil {
			messages.Error(install.Err.Error())
		}
		return
	}

	messages.Success("project dependencies installed")

	return
}
