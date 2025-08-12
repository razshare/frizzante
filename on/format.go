package on

import (
	"github.com/razshare/frizzante/cli/path"
	"github.com/razshare/frizzante/cli/state"
	"github.com/razshare/frizzante/tui/messages"
	"os"
	"os/exec"
)

func Format() {
	Touch()

	gofmt := exec.Command(path.Go("."), "fmt")
	gofmt.Env = append(os.Environ())
	gofmt.Stderr = os.Stderr
	gofmt.Stdout = os.Stdout
	gofmt.Stdin = os.Stdin
	err := gofmt.Run()
	if err != nil {
		messages.Fatal(err)
	}

	pretty := exec.Command(path.Bun(*state.App), "x", "prettier", "--write", ".")
	pretty.Dir = *state.App
	pretty.Env = append(os.Environ())
	pretty.Stderr = os.Stderr
	pretty.Stdout = os.Stdout
	pretty.Stdin = os.Stdin
	err = pretty.Run()
	if err != nil {
		messages.Fatal(err)
	}

	messages.Success("project formatted")
}
