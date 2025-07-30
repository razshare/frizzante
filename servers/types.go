package servers

import (
	"embed"
	"github.com/gorilla/websocket"
	"github.com/razshare/frizzante/archives"
	"github.com/razshare/frizzante/guards"
	"github.com/razshare/frizzante/notifiers"
	"github.com/razshare/frizzante/views"
	"net"
	"net/http"
	"time"
)

type Server struct {
	Address           string
	SecureAddress     string
	FormMaxMemory     int64
	HttpServer        *http.Server
	HttpMux           *http.ServeMux
	Connections       map[string]*net.Conn
	ReadTimeout       time.Duration
	WriteTimeout      time.Duration
	HeaderMaxMemory   int
	Certificate       string
	Key               string
	Notifier          *notifiers.Notifier
	Efs               embed.FS
	ViewConfiguration views.Configuration
	PublicRoot        string
	WsUpgrader        *websocket.Upgrader
	Guards            []guards.Guard
	SessionArchive    archives.Archive
}
