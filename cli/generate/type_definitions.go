package generate

import (
	"errors"
	"os"

	"github.com/razshare/frizzante/tui/messages"
)

func TypeDefinitions(options TypeDefinitionsOptions) (err error) {
	if !messages.Command(messages.CommandOptions{
		Environment: append(os.Environ(), "DEV=1"),
		Program:     options.Go,
		Args:        []string{"run", "-tags=no_servers,types", "."},
	}) {
		err = errors.New("could not generate type definitions")
		return
	}
	messages.Success("type definitions generated")
	return
}
