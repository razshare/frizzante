package actions

import (
	"errors"
	"os"

	"github.com/razshare/frizzante/tui/messages"
	"github.com/razshare/frizzante/tui/spinners"
)

func Update(options UpdateOptions) (err error) {
	if err = Touch(TouchOptions{}); err != nil {
		return
	}
	spin := spinners.New("updating go packages")
	go spinners.Start(spin)
	if !messages.Command(messages.CommandOptions{
		Environment: os.Environ(),
		Program:     options.Go,
		Args:        []string{"mod", "tidy"},
	}) {
		err = errors.New("could not update go packages")
		return
	}
	if !messages.Command(messages.CommandOptions{
		Environment: os.Environ(),
		Program:     options.Go,
		Args:        []string{"get", "-U", "./..."},
	}) {
		err = errors.New("could not update go packages")
		return
	}
	messages.Success("go packages updated")
	spinners.Stop(spin)
	spin = spinners.New("updating javascript packages")
	go spinners.Start(spin)
	if messages.Command(messages.CommandOptions{
		DirectoryName: "app",
		Environment:   os.Environ(),
		Program:       options.Bun,
		Args:          []string{"update"},
	}) {
		messages.Success("javascript packages updated")
	}
	spinners.Stop(spin)
	return
}
