package receive

import (
	"io"
	"main/lib/core/client"

	"main/lib/core/stack"
)

// Message reads the contents of the message and returns the value.
//
// Compatible with web sockets.
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
