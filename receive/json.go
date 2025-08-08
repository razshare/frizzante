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
func Json(c *client.Client, v any) {
	if c.Scope.WebSocket != nil {
		jsonError := c.Scope.WebSocket.ReadJSON(v)
		if jsonError != nil {
			c.Scope.Container.Config.ErrorLog.Println(jsonError, stack.Trace())
			return
		}
		return
	}

	data, readError := io.ReadAll(c.Request.Body)
	if readError != nil {
		c.Scope.Container.Config.ErrorLog.Println(readError, stack.Trace())
		return
	}

	jsonError := json.Unmarshal(data, v)
	if jsonError != nil {
		c.Scope.Container.Config.ErrorLog.Println(jsonError, stack.Trace())
		return
	}
}
