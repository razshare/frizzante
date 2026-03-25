package send

import (
	"github.com/gorilla/websocket"
	"github.com/razshare/frizzante/internal/project/lib/core/clients"
	"github.com/razshare/frizzante/internal/project/lib/core/stack"
)

// WsUpgrade upgrades to web sockets.
func WsUpgrade(client *clients.Client) {
	WsUpgradeWithUpgrader(client, websocket.Upgrader{
		ReadBufferSize:  10240, // 10KB
		WriteBufferSize: 10240, // 10KB
	})
}

// WsUpgradeWithUpgrader upgrades to web sockets.
func WsUpgradeWithUpgrader(client *clients.Client, upgrader websocket.Upgrader) {
	conn, err := upgrader.Upgrade(client.Writer, &client.Request, nil)
	if err != nil {
		client.Options.ErrorLog.Printf(
			"send.WsUpgradeWithUpgrader: failed to upgrade to WebSocket: %v\n%s",
			err,
			stack.Trace(),
		)
		return
	}

	defer func() {
		if cerr := conn.Close(); cerr != nil {
			client.Options.ErrorLog.Printf(
				"send.WsUpgradeWithUpgrader: failed to close WebSocket connection: %v\n%s",
				cerr,
				stack.Trace(),
			)
		}
	}()

	client.WebSocket = conn
	client.Locked = true
}
