package generate

import (
	"errors"
	"os"

	"github.com/razshare/frizzante/tui/messages"
)

func Definitions(options DefinitionsOptions) (err error) {
	if !messages.Command(".", append(os.Environ(), "DEV=1"), options.Go, "run", "-tags=dry,types", ".") {
		err = errors.New("could not generate definitions")
		return
	}

	messages.Success("definitions generated")

	return
}
