package send

import (
	"github.com/gorilla/websocket"
	"github.com/razshare/frizzante/conn"
	"github.com/razshare/frizzante/globals"
	"github.com/razshare/frizzante/stack"
)

// WsUpgrade upgrades to web sockets.
func WsUpgrade(c *conn.Conn) {
	WsUpgradeWithUpgrader(c, websocket.Upgrader{
		ReadBufferSize:  10 * globals.KB,
		WriteBufferSize: 10 * globals.KB,
	})
}

// WsUpgradeWithUpgrader upgrades to web sockets.
func WsUpgradeWithUpgrader(c *conn.Conn, u websocket.Upgrader) {
	webSocketConnection, upgradeError := u.Upgrade(c.Writer, c.Request, nil)
	if upgradeError != nil {
		c.Container.Config.ErrorLog.Println(upgradeError, stack.Trace())
		return
	}

	defer func(webSocketConnection *websocket.Conn) {
		closeError := c.WebSocket.Close()
		if closeError != nil {
			c.Container.Config.ErrorLog.Println(closeError, stack.Trace())
		}
	}(webSocketConnection)

	c.WebSocket = webSocketConnection
	c.Locked = true

	return
}
