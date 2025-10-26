package generate

import (
	"errors"
	"os"

	"github.com/razshare/frizzante/tui/messages"
)

func Queries(options QueriesOptions) (err error) {
	if !messages.Command("", append(os.Environ(), "DEV=1"), options.Go, "run", "-tags=dry,queries", ".") {
		err = errors.New("could not generate queries")
		return
	}

	messages.Success("queries generated")

	return
}
