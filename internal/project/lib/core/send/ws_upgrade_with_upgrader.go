package send

import (
	"github.com/gorilla/websocket"
	"github.com/razshare/frizzante/internal/project/lib/core/clients"
	"github.com/razshare/frizzante/internal/project/lib/core/logs"
	"github.com/razshare/frizzante/internal/project/lib/core/stack"
)

// WsUpgradeWithUpgrader upgrades to web sockets.
func WsUpgradeWithUpgrader(client *clients.Client, upgrader websocket.Upgrader) {
	conn, err := upgrader.Upgrade(client.Writer, &client.Request, nil)
	if err != nil {
		logs.Errorf(
			client,
			"send.WsUpgradeWithUpgrader: failed to upgrade to WebSocket: %v\n%s",
			err,
			stack.Trace(),
		)
		return
	}
	client.WebSocket = conn
	client.Locked = true
}
