package connections

import (
	"embed"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

type Connection struct {
	ErrorLog   *log.Logger
	InfoLog    *log.Logger
	WebSocket  *websocket.Conn
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
