package send

import (
	"github.com/gorilla/websocket"
	"github.com/razshare/frizzante/internal/project/lib/core/logs"
	"github.com/razshare/frizzante/internal/project/lib/core/scopes"
	"github.com/razshare/frizzante/internal/project/lib/core/stack"
)

// WsUpgradeWithUpgrader upgrades to web sockets.
func WsUpgradeWithUpgrader(http *scopes.Http, upgrader websocket.Upgrader) {
	conn, err := upgrader.Upgrade(http.Writer, &http.Request, nil)
	if err != nil {
		logs.Errorf(
			http,
			"send.WsUpgradeWithUpgrader: failed to upgrade to WebSocket: %v\n%s",
			err,
			stack.Trace(),
		)
		return
	}
	http.WebSocket = conn
	http.Locked = true
}
