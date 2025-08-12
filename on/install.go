package on

import (
	"github.com/razshare/frizzante/cli/path"
	"github.com/razshare/frizzante/cli/state"
	"github.com/razshare/frizzante/tui/messages"
	"github.com/razshare/frizzante/tui/spinner"
	"os"
	"os/exec"
)

func Install() {
	Touch()

	s := spinner.New("installing go dependencies")
	err := spinner.Start(s)
	if err != nil {
		spinner.Stop(s)
		messages.Fatal(err)
		return
	}
	tidy := exec.Command(path.Go("."), "mod", "tidy")
	tidy.Env = append(os.Environ())
	tidy.Stderr = os.Stderr
	tidy.Stdout = os.Stdout
	tidy.Stdin = os.Stdin
	err = tidy.Run()
	if err != nil {
		spinner.Stop(s)
		messages.Fatal(err)
	}
	spinner.Stop(s)

	ins := exec.Command(path.Bun(*state.App), "install")
	ins.Dir = *state.App
	ins.Env = append(os.Environ())
	ins.Stderr = os.Stderr
	ins.Stdout = os.Stdout
	ins.Stdin = os.Stdin
	err = ins.Run()
	if err != nil {
		messages.Fatal(err)
	}

	messages.Success("project dependencies installed")
}
