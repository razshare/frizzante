package on

import (
	"github.com/razshare/frizzante/cli/path"
	"github.com/razshare/frizzante/cli/state"
	"github.com/razshare/frizzante/tui/messages"
	"github.com/razshare/frizzante/tui/spinner"
	"os"
	"os/exec"
)

func Update() {
	Touch()
	s := spinner.New("updating go dependencies")
	err := spinner.Start(s)
	if err != nil {
		messages.Fatal(err)
		return
	}
	get := exec.Command(path.Go("."), "get", "-u", "./...")
	get.Env = append(os.Environ())
	get.Stderr = os.Stderr
	get.Stdout = os.Stdout
	get.Stdin = os.Stdin
	err = get.Run()
	if err != nil {
		spinner.Stop(s)
		messages.Fatal(err)
	}
	spinner.Stop(s)

	pretty := exec.Command(path.Bun(*state.App), "update")
	pretty.Dir = *state.App
	pretty.Env = append(os.Environ())
	pretty.Stderr = os.Stderr
	pretty.Stdout = os.Stdout
	pretty.Stdin = os.Stdin
	err = pretty.Run()
	if err != nil {
		messages.Fatal(err)
	}

	messages.Success("project dependencies updated")
}
