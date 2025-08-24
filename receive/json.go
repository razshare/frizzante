package receive

import (
	"encoding/json"
	"github.com/razshare/frizzante/client"
	"github.com/razshare/frizzante/stack"
	"io"
)

// Json reads the next JSON-encoded message from the
// c and stores it in the value pointed to by va.
//
// Compatible with web sockets.
func Json[T any](c *client.Client) T {
	var v T

	if c.WebSocket != nil {
		if err := c.WebSocket.ReadJSON(&v); err != nil {
			c.Config.ErrorLog.Println(err, stack.Trace())
			return v
		}
		return v
	}

	d, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.Config.ErrorLog.Println(err, stack.Trace())
		return v
	}

	if err = json.Unmarshal(d, &v); err != nil {
		c.Config.ErrorLog.Println(err, stack.Trace())
		return v
	}

	return v
}
