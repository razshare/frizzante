package send

import (
	"main/lib/core/client"

	"github.com/gorilla/websocket"
	"main/lib/core/stack"
)

// WsUpgrade upgrades to web sockets.
func WsUpgrade(client *client.Client) {
	WsUpgradeWithUpgrader(client, websocket.Upgrader{
		ReadBufferSize:  10240, // 10KB
		WriteBufferSize: 10240, // 10KB
	})
}

// WsUpgradeWithUpgrader upgrades to web sockets.
func WsUpgradeWithUpgrader(client *client.Client, upgrader websocket.Upgrader) {
	conn, err := upgrader.Upgrade(client.Writer, client.Request, nil)
	if err != nil {
		client.Config.ErrorLog.Println(err, stack.Trace())
		return
	}

	defer func() {
		if cerr := client.WebSocket.Close(); cerr != nil {
			client.Config.ErrorLog.Println(cerr, stack.Trace())
		}
	}()

	client.WebSocket = conn
	client.Locked = true

	return
}
