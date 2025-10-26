package paths

import (
	"os"

	"github.com/razshare/frizzante/tui/messages"
)

// Current gets the current working directory.
//
// When Current fails to retrieve the current working directory,
// it returns "." and logs the error.
func Current() string {
	if wd, err := os.Getwd(); err != nil {
		messages.Error(err)
		return "."
	} else {
		return wd
	}
}
