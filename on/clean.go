package on

import (
	"github.com/razshare/frizzante/cli/path"
	"github.com/razshare/frizzante/cli/state"
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
	err := clean.Run()
	if err != nil {
		messages.Fatal(err)
	}

	err = os.RemoveAll(filepath.Join(*state.App, "dist"))
	if err != nil {
		messages.Fatal(err)
	}

	err = os.RemoveAll(filepath.Join(*state.App, "node_modules"))
	if err != nil {
		messages.Fatal(err)
	}

	err = os.RemoveAll(filepath.Join(".gen", "tmp"))
	if err != nil {
		messages.Fatal(err)
	}

	err = os.RemoveAll(".vite")
	if err != nil {
		messages.Fatal(err)
	}

	Touch()

	messages.Success("project cleaned")
}
