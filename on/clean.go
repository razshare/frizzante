package on

import (
	"github.com/razshare/frizzante/cli/flags"
	"github.com/razshare/frizzante/cli/path"
	"github.com/razshare/frizzante/tui/messages"
	"os"
	"os/exec"
	"path/filepath"
)

func Clean() {
	clean := exec.Command(path.Go("."), "clean")
	clean.Env = append(os.Environ())
	clean.Stderr = os.Stderr
	clean.Stdout = os.Stdout
	clean.Stdin = os.Stdin
	runError := clean.Run()
	if runError != nil {
		messages.Fatal(runError)
	}

	removeError := os.RemoveAll(filepath.Join(*flags.App, "dist"))
	if removeError != nil {
		messages.Fatal(removeError)
	}

	removeError = os.RemoveAll(filepath.Join(*flags.App, "node_modules"))
	if removeError != nil {
		messages.Fatal(removeError)
	}

	removeError = os.RemoveAll(filepath.Join(".gen", "tmp"))
	if removeError != nil {
		messages.Fatal(removeError)
	}

	removeError = os.RemoveAll(".vite")
	if removeError != nil {
		messages.Fatal(removeError)
	}

	Touch()

	messages.Success("project cleaned")
}
