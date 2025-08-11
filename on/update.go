package on

import (
	"github.com/razshare/frizzante/cli/flags"
	"github.com/razshare/frizzante/cli/path"
	"github.com/razshare/frizzante/tui/messages"
	"os"
	"os/exec"
)

func Update() {
	Touch()

	get := exec.Command(path.Go("."), "get", "-u", "./...")
	get.Env = append(os.Environ())
	get.Stderr = os.Stderr
	get.Stdout = os.Stdout
	get.Stdin = os.Stdin
	getError := get.Run()
	if getError != nil {
		messages.Fatal(getError)
	}

	prettier := exec.Command(path.Bun(*flags.App), "update")
	prettier.Dir = *flags.App
	prettier.Env = append(os.Environ())
	prettier.Stderr = os.Stderr
	prettier.Stdout = os.Stdout
	prettier.Stdin = os.Stdin
	prettierError := prettier.Run()
	if prettierError != nil {
		messages.Fatal(prettierError)
	}

	messages.Success("project dependencies updated")
}
