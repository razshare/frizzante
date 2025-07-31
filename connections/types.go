package connections

import (
	"embed"
	"github.com/gorilla/websocket"
	"net/http"
)

type Connection struct {
	WebSocket  *websocket.Conn
	Http       *http.Server
	Request    *http.Request
	Writer     http.ResponseWriter
	Efs        embed.FS
	Status     int
	Locked     bool
	EventId    int64
	EventName  string
	SessionId  string
	PublicRoot string
	AppRoot    string
	ServerJs   string
	IndexHtml  string
}
