package actions

import (
	"errors"

	"github.com/razshare/frizzante/tui/messages"
)

func Snapshot(options SnapshotOptions) (err error) {
	if !messages.Command(messages.CommandOptions{Program: options.Program}) {
		err = errors.New("could not create snapshot")
		return
	}

	return
}
