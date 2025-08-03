package connections

import (
	"embed"
	"github.com/gorilla/websocket"
	"github.com/razshare/frizzante/archives"
	"github.com/razshare/frizzante/views"
	"log"
	"net/http"
)

type Connection struct {
	ErrorLog       *log.Logger
	InfoLog        *log.Logger
	WebSocket      *websocket.Conn
	Request        *http.Request
	Writer         http.ResponseWriter
	SessionArchive archives.Archive
	ViewContainer  *views.Container
	Efs            embed.FS
	Status         int
	Locked         bool
	EventId        int64
	EventName      string
	SessionId      string
	PublicRoot     string
}
