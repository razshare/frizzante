package connections

import (
	"embed"
	"github.com/gorilla/websocket"
	"github.com/razshare/frizzante/apps"
	"github.com/razshare/frizzante/archives"
	"log"
	"net/http"
)

type Connection struct {
	Status           int
	Locked           bool
	EventId          int64
	EventName        string
	SessionId        string
	PublicRoot       string
	Efs              embed.FS
	App              *apps.App
	ErrorLog         *log.Logger
	InfoLog          *log.Logger
	Request          *http.Request
	WebSocket        *websocket.Conn
	SessionArchive   archives.Archive
	AppConfiguration apps.Configuration
	Writer           http.ResponseWriter
}
