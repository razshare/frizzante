package clients

import (
	"net/http"

	"github.com/gorilla/websocket"
)

type Client struct {
	Request   http.Request
	Options   Options
	Writer    http.ResponseWriter
	SessionId string
	EventName string
	EventId   int64
	Status    int
	WebSocket *websocket.Conn
	Locked    bool
	Parsed    bool
}
