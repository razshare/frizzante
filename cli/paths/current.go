package paths

import (
	"os"

	"github.com/razshare/frizzante/v2/tui/messages"
)

// Current gets the current working directory.
//
// When Current fails to retrieve the current working directory,
// it returns "." and logs the error.
func Current() (wd string) {
	var err error
	if wd, err = os.Getwd(); err != nil {
		messages.Error(err)
		return "."
	}
	return
}
