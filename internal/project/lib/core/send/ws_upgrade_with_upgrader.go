package send

import (
	"net/http"

	"github.com/gorilla/websocket"
)

// WsUpgradeWithUpgrader upgrades to web sockets.
func WsUpgradeWithUpgrader(
	writer *http.ResponseWriter,
	request *http.Request,
	upgrader websocket.Upgrader,
) (err error) {
	var webSocketConnection *websocket.Conn
	if webSocketConnection, err = upgrader.Upgrade(*writer, request, nil); err != nil {
		return
	}
	converted := http.ResponseWriter(&WsUpgradeResponseWriter{
		ResponseWriter:      *writer,
		WebSocketConnection: webSocketConnection,
	})
	request.Body = &WsUpgradeBodyReaderCloser{
		WebSocketConnection: webSocketConnection,
	}
	writer = &converted
	return
}
