package actions

import (
	"errors"
	"os"

	"github.com/razshare/frizzante/tui/messages"
)

func Test(options TestOptions) (err error) {
	if !messages.Command(messages.CommandOptions{
		Env:  os.Environ(),
		Name: options.Go,
		Args: []string{"test", "./..."},
	}) {
		err = errors.New("tests failed")
		return
	}
	return
}
