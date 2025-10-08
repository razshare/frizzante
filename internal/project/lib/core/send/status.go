package send

import (
	"github.com/razshare/frizzante/internal/project/lib/core/clients"
	"github.com/razshare/frizzante/internal/project/lib/core/stacks"
)

// Status sets the status code.
//
// This will lock the status, which makes it
// so that the next time you invoke this
// function it will fail with an error.
//
// All errors are sent to the server notifier.
func Status(client *clients.Client, status int) {
	if client.Locked {
		client.Config.ErrorLog.Println("status is locked", stacks.Trace())
		return
	}

	client.Status = status
}
