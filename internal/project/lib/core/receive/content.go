package receive

import (
	"io"

	"github.com/razshare/frizzante/internal/project/lib/core/client"
	"github.com/razshare/frizzante/internal/project/lib/core/stack"
)

// Message reads the contents of the message and returns the value.
//
// Compatible with web sockets and server sent events.
func Message(client *client.Client) string {
	if client.WebSocket != nil {
		_, data, err := client.WebSocket.ReadMessage()
		if err != nil {
			client.Config.ErrorLog.Println(err, stack.Trace())
			return ""
		}
		return string(data)
	}

	data, err := io.ReadAll(client.Request.Body)
	if err != nil {
		client.Config.ErrorLog.Println(err, stack.Trace())
		return ""
	}
	return string(data)
}
