package send

import (
	"github.com/gorilla/websocket"
	"github.com/razshare/frizzante/internal/project/lib/core/clients"
	"github.com/razshare/frizzante/internal/project/lib/core/stack"
)

// Content sends binary safe content.
//
// Compatible with web sockets and server sent events.
func Content(client *clients.Client, data []byte) {
	if !client.Locked {
		client.Writer.WriteHeader(client.Status)
		client.Locked = true
	}
	if client.WebSocket != nil {
		if err := client.WebSocket.WriteMessage(websocket.TextMessage, data); err != nil {
			client.Options.ErrorLog.Printf(
				"send.Content: failed to write websocket message: %v\n%s",
				err,
				stack.Trace(),
			)
		}
		return
	}
	if client.EventName != "" {
		EventContent(client, data)
		return
	}
	if _, err := client.Writer.Write(data); err != nil {
		client.Options.ErrorLog.Printf(
			"send.Content: failed to write response content: %v\n%s",
			err,
			stack.Trace(),
		)
	}
}
