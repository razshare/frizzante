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

	if !messages.Command(messages.CommandOptions{
		Env:  os.Environ(),
		Name: options.Go,
		Args: []string{"test", "./..."},
	}) {
		spinners.Stop(spin)
		err = errors.New("tests failed")
		return
	}
	spinners.Stop(spin)
	return
}
