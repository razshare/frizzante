package send

import (
	"github.com/gorilla/websocket"
	"github.com/razshare/frizzante/client"
	"github.com/razshare/frizzante/globals"
	"github.com/razshare/frizzante/stack"
)

// WsUpgrade upgrades to web sockets.
func WsUpgrade(c *client.Client) {
	WsUpgradeWithUpgrader(c, websocket.Upgrader{
		ReadBufferSize:  10 * globals.KB,
		WriteBufferSize: 10 * globals.KB,
	})
}

// WsUpgradeWithUpgrader upgrades to web sockets.
func WsUpgradeWithUpgrader(c *client.Client, u websocket.Upgrader) {
	webSocketConnection, upgradeError := u.Upgrade(c.Writer, c.Request, nil)
	if upgradeError != nil {
		c.Scope.Container.Config.ErrorLog.Println(upgradeError, stack.Trace())
		return
	}

	defer func(webSocketConnection *websocket.Conn) {
		closeError := c.Scope.WebSocket.Close()
		if closeError != nil {
			c.Scope.Container.Config.ErrorLog.Println(closeError, stack.Trace())
		}
	}(webSocketConnection)

	c.Scope.WebSocket = webSocketConnection
	c.Scope.Locked = true

	return
}
