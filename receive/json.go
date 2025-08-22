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
func Json[T any](client *client.Client) T {
	var val T

	if client.WebSocket != nil {
		if err := client.WebSocket.ReadJSON(&val); err != nil {
			client.Config.ErrorLog.Println(err, stack.Trace())
			return val
		}
		return val
	}

	data, err := io.ReadAll(client.Request.Body)
	if err != nil {
		client.Config.ErrorLog.Println(err, stack.Trace())
		return val
	}

	if err = json.Unmarshal(data, &val); err != nil {
		client.Config.ErrorLog.Println(err, stack.Trace())
		return val
	}

	return val
}
