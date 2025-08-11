package on

import (
	"github.com/razshare/frizzante/cli/path"
	"github.com/razshare/frizzante/tui/messages"
	"os"
	"os/exec"
)

func SqlcGenerate() {
	sqlcGenerate := exec.Command(path.Sqlc("."), "generate")
	sqlcGenerate.Env = append(os.Environ())
	sqlcGenerate.Stderr = os.Stderr
	sqlcGenerate.Stdout = os.Stdout
	sqlcGenerate.Stdin = os.Stdin
	svelteCheckError := sqlcGenerate.Run()
	if svelteCheckError != nil {
		messages.Fatal(svelteCheckError)
	}
	messages.Success("files generated")
}
