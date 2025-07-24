package connections

import (
	"embed"
	"github.com/gorilla/websocket"
	"github.com/razshare/frizzante/archives"
	"github.com/razshare/frizzante/notifiers"
	"net/http"
)

type Connection struct {
	PublicRoot     string
	Notifier       *notifiers.Notifier
	Efs            embed.FS
	ViewServer     string
	ViewIndex      string
	ViewRoot       string
	Request        *http.Request
	Writer         http.ResponseWriter
	Locked         bool
	Status         int
	Header         http.Header
	WebSocket      *websocket.Conn
	EventName      string
	EventId        int64
	SessionId      string
	SessionArchive archives.Archive
}
