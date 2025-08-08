package receive

import (
	"github.com/razshare/frizzante/client"
	"github.com/razshare/frizzante/stack"
	"io"
)

// Message reads the contents of the message and returns the value.
//
// Compatible with web sockets.
func Message(c *client.Client) string {
	if c.Scope.WebSocket != nil {
		_, data, readError := c.Scope.WebSocket.ReadMessage()
		if readError != nil {
			c.Scope.Container.Config.ErrorLog.Println(readError, stack.Trace())
			return ""
		}
		return string(data)
	}

	data, readError := io.ReadAll(c.Request.Body)
	if readError != nil {
		c.Scope.Container.Config.ErrorLog.Println(readError, stack.Trace())
		return ""
	}
	return string(data)
}
