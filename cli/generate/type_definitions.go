package generate

import (
	"errors"
	"os"

	"github.com/razshare/frizzante/tui/messages"
)

func TypeDefinitions(options TypeDefinitionsOptions) (err error) {
	if !messages.Command("", append(os.Environ(), "DEV=1"), options.Go, "run", "-tags=dry,types", ".") {
		err = errors.New("could not generate types")
		return
	}

	messages.Success("types generated")

	return
}
