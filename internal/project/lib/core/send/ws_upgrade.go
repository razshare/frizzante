package send

import (
	"github.com/gorilla/websocket"
	"github.com/razshare/frizzante/internal/project/lib/core/clients"
)

// WsUpgrade upgrades to web sockets.
func WsUpgrade(client *clients.Client) {
	WsUpgradeWithUpgrader(client, websocket.Upgrader{
		ReadBufferSize:  10240, // 10KB
		WriteBufferSize: 10240, // 10KB
	})
}
