package actions

import (
	"errors"
	"os"

	"github.com/razshare/frizzante/v2/tui/messages"
	"github.com/razshare/frizzante/v2/tui/spinners"
)

func Format(options FormatOptions) (err error) {
	spin := spinners.New("formatting code")
	go spinners.Start(spin)
	defer spinners.Stop(spin)
	if !messages.Command(messages.CommandOptions{
		Program:     options.Go,
		Environment: os.Environ(),
		Args:        []string{"fmt", "./..."},
	}) {
		err = errors.New("could not format go code")
		return
	}
	if !messages.Command(messages.CommandOptions{
		DirectoryName: "app",
		Environment:   os.Environ(),
		Program:       options.Bun,
		Args:          []string{"x", "prettier", "--write", "."},
	}) {
		err = errors.New("could not format js code")
		return
	}
	messages.Success("project formatted")
	return
}
