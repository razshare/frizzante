package send

import (
	"github.com/gorilla/websocket"
	"github.com/razshare/frizzante/internal/project/lib/core/clients"
	"github.com/razshare/frizzante/internal/project/lib/core/stacks"
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
	conn, err := upgrader.Upgrade(client.Writer, client.Request, nil)
	if err != nil {
		client.Config.ErrorLog.Println(err, stacks.Trace())
		return
	}

	defer func() {
		if cerr := client.WebSocket.Close(); cerr != nil {
			client.Config.ErrorLog.Println(cerr, stacks.Trace())
		}
	}()

	client.WebSocket = conn
	client.Locked = true
}
