package action

import (
	"github.com/razshare/frizzante/tui/messages"
	"os"
	"os/exec"
	"path/filepath"
)

func CleanProject(o CleanProjectOptions) error {
	clean := exec.Command(o.Go, "clean")
	clean.Env = append(os.Environ())
	clean.Stderr = os.Stderr
	clean.Stdout = os.Stdout
	clean.Stdin = os.Stdin
	err := clean.Run()
	if err != nil {
		return err
	}

	err = os.RemoveAll(filepath.Join(o.App, "dist"))
	if err != nil {
		return err
	}

	err = os.RemoveAll(filepath.Join(o.App, "node_modules"))
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

	err = Touch(TouchOptions{App: o.App})
	if err != nil {
		return err
	}

	messages.Success("project cleaned")

	return nil
}
