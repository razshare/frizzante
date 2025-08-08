package conn

import (
	"github.com/gorilla/websocket"
	"github.com/razshare/frizzante/container"
	"net/http"
)

type Conn struct {
	Container *container.Container
	Writer    http.ResponseWriter
	WebSocket *websocket.Conn
	Request   *http.Request
	EventName string
	EventId   int64
	Locked    bool
	Status    int
}
