package connections

import (
	"embed"
	"github.com/gorilla/websocket"
	"github.com/razshare/frizzante/apps"
	"log"
	"net/http"
)

type Connection struct {
	Writer     http.ResponseWriter
	AppConfig  *apps.Config
	WebSocket  *websocket.Conn
	Request    *http.Request
	ErrorLog   *log.Logger
	App        *apps.App
	Efs        embed.FS
	PublicRoot string
	EventName  string
	EventId    int64
	Locked     bool
	Status     int
}
