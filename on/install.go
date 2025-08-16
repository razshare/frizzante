package on

import (
	"github.com/razshare/frizzante/cli/path"
	"github.com/razshare/frizzante/cli/state"
	"github.com/razshare/frizzante/tui/messages"
	"github.com/razshare/frizzante/tui/spinner"
	"os"
	"os/exec"
	"path/filepath"
)

func Install(base string) error {
	err := Touch()
	if err != nil {
		return err
	}

	s := spinner.New("installing go dependencies")

	gobin, err := path.Go(base)
	if err != nil {
		return err
	}

	go spinner.Start(s)
	tidy := exec.Command(gobin, "mod", "tidy")
	tidy.Dir = base
	tidy.Env = append(os.Environ())
	tidy.Stderr = os.Stderr
	tidy.Stdout = os.Stdout
	tidy.Stdin = os.Stdin
	err = tidy.Run()
	spinner.Stop(s)

	if err != nil {
		return err
	}

	gobin, err = path.Bun(filepath.Join(base, *state.App))
	if err != nil {
		return err
	}

	ins := exec.Command(gobin, "install")
	ins.Dir = filepath.Join(base, *state.App)
	ins.Env = append(os.Environ())
	ins.Stderr = os.Stderr
	ins.Stdout = os.Stdout
	ins.Stdin = os.Stdin
	err = ins.Run()
	if err != nil {
		return err
	}

	messages.Success("project dependencies installed")

	return nil
}
