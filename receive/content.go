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
	if c.WebSocket != nil {
		_, d, err := c.WebSocket.ReadMessage()
		if err != nil {
			c.Config.ErrorLog.Println(err, stack.Trace())
			return ""
		}
		return string(d)
	}

	d, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.Config.ErrorLog.Println(err, stack.Trace())
		return ""
	}
	return string(d)
}
