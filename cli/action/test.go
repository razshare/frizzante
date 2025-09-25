package action

import (
	"errors"
	"os"

	"github.com/razshare/frizzante/tui/messages"
	"github.com/razshare/frizzante/tui/spinner"
)

func Test(options TestOptions) (err error) {
	spin := spinner.New("testing code")
	go spinner.Start(spin)
	defer spinner.Stop(spin)

	if !messages.Command(".", os.Environ(), options.Go, "test", "./...") {
		err = errors.New("tests failed")
	}
	return
}
