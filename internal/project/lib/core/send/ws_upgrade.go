package send

import (
	"github.com/gorilla/websocket"
	"github.com/razshare/frizzante/internal/project/lib/core/scopes"
)

// WsUpgrade upgrades to web sockets.
func WsUpgrade(http *scopes.Http) {
	WsUpgradeWithUpgrader(http, websocket.Upgrader{
		ReadBufferSize:  10240, // 10KB
		WriteBufferSize: 10240, // 10KB
	})
}
