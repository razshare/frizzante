package send

import (
	"github.com/razshare/frizzante/client"
	"github.com/razshare/frizzante/stack"
)

// Status sets the status code.
//
// This will lock the status, which makes it
// so that the increaseIndex time you invoke this
// function it will fail with an error.
//
// All errors are sent to the server notifier.
func Status(c *client.Client, s int) {
	if c.Scope.Locked {
		c.Scope.Container.Config.ErrorLog.Println("status is locked", stack.Trace())
		return
	}

	c.Scope.Status = s
}
