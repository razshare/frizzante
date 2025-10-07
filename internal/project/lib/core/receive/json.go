package receive

import (
	"encoding/json"
	"io"

	_client "github.com/razshare/frizzante/internal/project/lib/core/client"
	"github.com/razshare/frizzante/internal/project/lib/core/stack"
)

// Json reads the next JSON-encoded message from the
// c and stores it in the value pointed to by value.
//
// Compatible with web sockets and server sent events.
func Json(client *_client.Client, value any) bool {
	if client.WebSocket != nil {
		if err := client.WebSocket.ReadJSON(&value); err != nil {
			client.Config.ErrorLog.Println(err, stack.Trace())
			return false
		}
		return true
	}

	data, err := io.ReadAll(client.Request.Body)
	if err != nil {
		client.Config.ErrorLog.Println(err, stack.Trace())
		return false
	}

	if err = json.Unmarshal(data, &value); err != nil {
		client.Config.ErrorLog.Println(err, stack.Trace())
	}

	return true
}
