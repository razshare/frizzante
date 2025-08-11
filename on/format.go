package on

import (
	"github.com/razshare/frizzante/cli/flags"
	"github.com/razshare/frizzante/cli/path"
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
	gofmtError := gofmt.Run()
	if gofmtError != nil {
		messages.Fatal(gofmtError)
	}

	prettier := exec.Command(path.Bun(*flags.App), "x", "prettier", "--write", ".")
	prettier.Dir = *flags.App
	prettier.Env = append(os.Environ())
	prettier.Stderr = os.Stderr
	prettier.Stdout = os.Stdout
	prettier.Stdin = os.Stdin
	prettierError := prettier.Run()
	if prettierError != nil {
		messages.Fatal(prettierError)
	}

	messages.Success("project formatted")
}
