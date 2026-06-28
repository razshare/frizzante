package send

import (
	"net/http"

	"github.com/gorilla/websocket"
)

// WsUpgrade upgrades to web sockets.
func WsUpgrade(writer *http.ResponseWriter, request *http.Request) error {
	return WsUpgradeWithUpgrader(writer, request, websocket.Upgrader{
		ReadBufferSize:  10240, // 10KB
		WriteBufferSize: 10240, // 10KB
	})
}
