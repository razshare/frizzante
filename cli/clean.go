package cli

import (
	"os"
	"os/exec"
	"path/filepath"
)

func OnClean() {
	clean := exec.Command(Go("."), "clean")
	clean.Env = append(os.Environ())
	clean.Stderr = os.Stderr
	clean.Stdout = os.Stdout
	clean.Stdin = os.Stdin
	runError := clean.Run()
	if runError != nil {
		Fatal(runError)
	}

	removeError := os.RemoveAll(filepath.Join("app", "dist"))
	if removeError != nil {
		Fatal(removeError)
	}

	removeError = os.RemoveAll(filepath.Join("app", "node_modules"))
	if removeError != nil {
		Fatal(removeError)
	}

	removeError = os.RemoveAll(filepath.Join(".gen", "tmp"))
	if removeError != nil {
		Fatal(removeError)
	}

	removeError = os.RemoveAll(".vite")
	if removeError != nil {
		Fatal(removeError)
	}

	OnTouch()

	Success("project cleaned")
}
