package servers

import (
	"embed"
	"github.com/gorilla/websocket"
	"github.com/razshare/frizzante/archives"
	"github.com/razshare/frizzante/guards"
	"github.com/razshare/frizzante/notifiers"
	"net"
	"net/http"
	"time"
)

type Server struct {
	Address         string
	SecureAddress   string
	FormMaxMemory   int64
	HttpServer      *http.Server
	HttpMux         *http.ServeMux
	Connections     map[string]*net.Conn
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	HeaderMaxMemory int
	Certificate     string
	Key             string
	Notifier        *notifiers.Notifier
	Efs             embed.FS
	WsUpgrader      *websocket.Upgrader
	Guards          []guards.Guard
	SessionArchive  archives.Archive
	PublicRoot      string
	AppRoot         string
	ServerJs        string
	IndexHtml       string
	Dotenv          string
}
