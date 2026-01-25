package actions

import (
	"errors"

	"github.com/razshare/frizzante/tui/messages"
	"github.com/razshare/frizzante/tui/spinners"
)

func Format(options FormatOptions) (err error) {
	spin := spinners.New("formatting code")
	go spinners.Start(spin)
	defer spinners.Stop(spin)
	if err = Touch(TouchOptions{}); err != nil {
		return
	}
	if !messages.Command(messages.CommandOptions{
		Program: options.Go,
		Args:    []string{"fmt", "./..."},
	}) {
		err = errors.New("could not format go code")
		return
	}
	if !messages.Command(messages.CommandOptions{
		DirectoryName: "app",
		Program:       options.Bun,
		Args:          []string{"x", "prettier", "--write", "."},
	}) {
		err = errors.New("could not format js code")
		return
	}
	messages.Success("project formatted")
	return
}
