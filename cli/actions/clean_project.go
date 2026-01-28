package actions

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/razshare/frizzante/tui/messages"
)

func CleanProject(options CleanProjectOptions) (err error) {
	if !messages.Command(messages.CommandOptions{Program: options.Go, Args: []string{"clean"}}) {
		err = errors.New("could not clean Go project")
		return
	}
	if err = os.RemoveAll(".gen"); err != nil {
		return
	}
	if err = os.RemoveAll(filepath.Join("dist")); err != nil {
		return
	}
	if err = os.RemoveAll(filepath.Join("node_modules")); err != nil {
		return
	}
	if err = os.RemoveAll(filepath.Join(".vite")); err != nil {
		return
	}
	if err = os.RemoveAll(filepath.Join("lib", "core", "ssr", "app")); err != nil {
		return
	}
	if err = os.MkdirAll(filepath.Join("app", "dist"), os.ModePerm); err != nil {
		return
	}
	messages.Success("project cleaned")
	return
}
