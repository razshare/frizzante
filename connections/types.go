package connections

import (
	"embed"
	"github.com/gorilla/websocket"
	"github.com/razshare/frizzante/archives"
	"github.com/razshare/frizzante/containers"
	"log"
	"net/http"
)

type Connection struct {
	Status         int
	Locked         bool
	EventId        int64
	EventName      string
	SessionId      string
	PublicRoot     string
	Efs            embed.FS
	ErrorLog       *log.Logger
	InfoLog        *log.Logger
	WebSocket      *websocket.Conn
	SessionArchive archives.Archive
	Request        *http.Request
	Writer         http.ResponseWriter
	ViewContainer  *containers.ViewContainer
}
