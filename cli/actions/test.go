package actions

import (
	"errors"
	"os"

	"github.com/razshare/frizzante/tui/messages"
	"github.com/razshare/frizzante/tui/spinners"
)

func Test(options TestOptions) (err error) {
	spin := spinners.New("testing code")
	go spinners.Start(spin)
	defer spinners.Stop(spin)

	if !messages.Command("", os.Environ(), options.Go, "test", "./...") {
		err = errors.New("tests failed")
	}
	return
}
