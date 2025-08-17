package on

import (
	"github.com/razshare/frizzante/tui/messages"
	"os"
	"os/exec"
	"path/filepath"
)

func Clean(app string, gobin string) error {
	clean := exec.Command(gobin, "clean")
	clean.Env = append(os.Environ())
	clean.Stderr = os.Stderr
	clean.Stdout = os.Stdout
	clean.Stdin = os.Stdin
	err := clean.Run()
	if err != nil {
		return err
	}

	err = os.RemoveAll(filepath.Join(app, "dist"))
	if err != nil {
		return err
	}

	err = os.RemoveAll(filepath.Join(app, "node_modules"))
	if err != nil {
		return err
	}

	err = os.RemoveAll(filepath.Join(".gen", "tmp"))
	if err != nil {
		return err
	}

	err = os.RemoveAll(filepath.Join(".vite"))
	if err != nil {
		return err
	}

	err = Touch(app)
	if err != nil {
		return err
	}

	messages.Success("project cleaned")

	return nil
}
