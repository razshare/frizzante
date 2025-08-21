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
		jsonError := c.WebSocket.ReadJSON(&v)
		if jsonError != nil {
			c.Config.ErrorLog.Println(jsonError, stack.Trace())
			return v
		}
		return v
	}

	data, readError := io.ReadAll(c.Request.Body)
	if readError != nil {
		c.Config.ErrorLog.Println(readError, stack.Trace())
		return v
	}

	jsonError := json.Unmarshal(data, &v)
	if jsonError != nil {
		c.Config.ErrorLog.Println(jsonError, stack.Trace())
		return v
	}

	return v
}
