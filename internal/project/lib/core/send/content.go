package send

import (
	"github.com/gorilla/websocket"
	"github.com/razshare/frizzante/internal/project/lib/core/logs"
	"github.com/razshare/frizzante/internal/project/lib/core/scopes"
	"github.com/razshare/frizzante/internal/project/lib/core/stack"
)

// Content sends binary safe content.
//
// Compatible with web sockets and server sent events.
func Content(http *scopes.Http, data []byte) {
	if !http.Locked {
		http.Writer.WriteHeader(http.Status)
		http.Locked = true
	}
	if http.WebSocket != nil {
		if err := http.WebSocket.WriteMessage(websocket.TextMessage, data); err != nil {
			logs.Errorf(
				http,
				"send.Content: failed to write websocket message: %v\n%s",
				err,
				stack.Trace(),
			)
		}
		return
	}
	if http.EventName != "" {
		EventContent(http, data)
		return
	}
	if _, err := http.Writer.Write(data); err != nil {
		logs.Errorf(
			http,
			"send.Content: failed to write response content: %v\n%s",
			err,
			stack.Trace(),
		)
	}
}
