package actions

import (
	"errors"
	"os"

	"github.com/razshare/frizzante/tui/messages"
	"github.com/razshare/frizzante/tui/spinners"
)

func Install(options InstallOptions) (err error) {
	if err = Touch(TouchOptions{}); err != nil {
		return
	}

	spin := spinners.New("installing go packages")
	go spinners.Start(spin)
	if !messages.Command(messages.CommandOptions{
		Environment: os.Environ(),
		Program:     options.Go,
		Args:        []string{"mod", "tidy"},
	}) {
		err = errors.New("could not install go packages")
		return
	}
	if !messages.Command(messages.CommandOptions{
		Environment: os.Environ(),
		Program:     options.Go,
		Args:        []string{"get", "./..."},
	}) {
		err = errors.New("could not install go packages")
		return
	}
	messages.Success("go packages installed")
	spinners.Stop(spin)

	spin = spinners.New("installing javascript packages")
	go spinners.Start(spin)
	if messages.Command(messages.CommandOptions{
		DirectoryName: "app",
		Environment:   os.Environ(),
		Program:       options.Bun,
		Args:          []string{"install"},
	}) {
		messages.Success("javascript packages installed")
	}
	spinners.Stop(spin)

	return
}
