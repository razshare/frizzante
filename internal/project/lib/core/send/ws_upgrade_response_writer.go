package send

import (
	"net/http"

	"github.com/gorilla/websocket"
)

type WsUpgradeResponseWriter struct {
	http.ResponseWriter
	WebSocketConnection *websocket.Conn
}

func (writer *WsUpgradeResponseWriter) Write(data []byte) (written int, err error) {
	written = len(data)
	err = writer.WebSocketConnection.WriteMessage(websocket.TextMessage, data)
	return
}
