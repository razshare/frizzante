package on

import (
	"github.com/razshare/frizzante/cli/flags"
	"github.com/razshare/frizzante/cli/path"
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

	install := exec.Command(path.Bun(*flags.App), "install")
	install.Dir = *flags.App
	install.Env = append(os.Environ())
	install.Stderr = os.Stderr
	install.Stdout = os.Stdout
	install.Stdin = os.Stdin
	err = install.Run()
	if err != nil {
		messages.Fatal(err)
	}

	messages.Success("project dependencies installed")
}
