package receive

import (
	"io"

	"github.com/razshare/frizzante/internal/project/lib/core/clients"
	"github.com/razshare/frizzante/internal/project/lib/core/logs"
	"github.com/razshare/frizzante/internal/project/lib/core/stack"
)

// Message reads the contents of the message and returns the value.
//
// Compatible with web sockets.
func Message(client *clients.Client) string {
	if client.WebSocket != nil {
		_, data, err := client.WebSocket.ReadMessage()
		if err != nil {
			logs.Errorf(
				client,
				"receive.Message: failed to read WebSocket message: %v\n%s",
				err,
				stack.Trace(),
			)
			return ""
		}
		return string(data)
	}
	data, err := io.ReadAll(client.Request.Body)
	if err != nil {
		logs.Errorf(
			client,
			"receive.Message: failed to read request body: %v\n%s",
			err,
			stack.Trace(),
		)
		return ""
	}
	return string(data)
}
