package receive

import (
	"encoding/json"
	"io"
	"main/lib/core/client"

	"github.com/razshare/frizzante/stack"
)

// Json reads the next JSON-encoded message from the
// c and stores it in the value pointed to by value.
//
// Compatible with web sockets.
func Json(client *client.Client, value any) {
	if client.WebSocket != nil {
		if err := client.WebSocket.ReadJSON(&value); err != nil {
			client.Config.ErrorLog.Println(err, stack.Trace())
			return
		}
		return
	}

	data, err := io.ReadAll(client.Request.Body)
	if err != nil {
		client.Config.ErrorLog.Println(err, stack.Trace())
		return
	}

	if err = json.Unmarshal(data, &value); err != nil {
		client.Config.ErrorLog.Println(err, stack.Trace())
	}
}
