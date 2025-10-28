package generate

import (
	"errors"
	"os"

	"github.com/razshare/frizzante/tui/messages"
)

func TypeDefinitions(options TypeDefinitionsOptions) (err error) {
	if !messages.Command(messages.CommandOptions{
		Env:  append(os.Environ(), "DEV=1"),
		Name: options.Go,
		Args: []string{"run", "-tags=dry,types", "."},
	}) {
		err = errors.New("could not generate types")
		return
	}
	messages.Success("types generated")
	return
}
