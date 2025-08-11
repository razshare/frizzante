package on

import (
	"github.com/razshare/frizzante/cli/flags"
	"github.com/razshare/frizzante/cli/path"
	"github.com/razshare/frizzante/tui/messages"
	"os"
	"os/exec"
)

func Install() {
	Touch()

	tidy := exec.Command(path.Go("."), "mod", "tidy")
	tidy.Env = append(os.Environ())
	tidy.Stderr = os.Stderr
	tidy.Stdout = os.Stdout
	tidy.Stdin = os.Stdin
	tidyError := tidy.Run()
	if tidyError != nil {
		messages.Fatal(tidyError)
	}

	install := exec.Command(path.Bun(*flags.App), "install")
	install.Dir = *flags.App
	install.Env = append(os.Environ())
	install.Stderr = os.Stderr
	install.Stdout = os.Stdout
	install.Stdin = os.Stdin
	installError := install.Run()
	if installError != nil {
		messages.Fatal(installError)
	}

	messages.Success("project dependencies installed")
}
