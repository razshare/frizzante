package actions

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/razshare/frizzante/v2/tui/messages"
)

func CleanProject(options CleanProjectOptions) (err error) {
	if !messages.Command(messages.CommandOptions{Program: options.Go, Args: []string{"clean"}}) {
		err = errors.New("could not clean Go project")
		return
	}
	if err = os.RemoveAll(".gen"); err != nil {
		return
	}
	if err = os.RemoveAll(filepath.Join("app", "dist")); err != nil {
		return
	}
	if err = os.RemoveAll(filepath.Join("app", ".vite")); err != nil {
		return
	}
	if err = os.RemoveAll(filepath.Join("app", "node_modules")); err != nil {
		return
	}
	if err = os.RemoveAll("cover.html"); err != nil {
		return
	}
	if err = os.RemoveAll("cover.out"); err != nil {
		return
	}
	messages.Success("project cleaned")
	return
}
