package receive

import (
	"io"

	"github.com/razshare/frizzante/internal/project/lib/core/clients"
	"github.com/razshare/frizzante/internal/project/lib/core/stacks"
)

// Message reads the contents of the message and returns the value.
//
// Compatible with web sockets and server sent events.
func Message(client *clients.Client) string {
	if client.WebSocket != nil {
		_, data, err := client.WebSocket.ReadMessage()
		if err != nil {
			client.Config.ErrorLog.Println(err, stacks.Trace())
			return ""
		}
		return string(data)
	}

	data, err := io.ReadAll(client.Request.Body)
	if err != nil {
		client.Config.ErrorLog.Println(err, stacks.Trace())
		return ""
	}
	return string(data)
}
