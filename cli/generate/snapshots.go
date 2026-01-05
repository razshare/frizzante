package generate

import (
	"errors"
	"os"

	"github.com/razshare/frizzante/tui/messages"
)

func Snapshots(options SnapshotsOptions) (err error) {
	if !messages.Command(messages.CommandOptions{
		Environment: append(os.Environ(), "DEV=1"),
		Program:     options.Go,
		Args:        []string{"run", "-tags=snapshots", "."},
	}) {
		err = errors.New("could not generate snapshots")
		return
	}
	messages.Success("snapshots generated")
	return
}
