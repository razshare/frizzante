package libsrv

import (
	"embed"
	"github.com/gorilla/websocket"
	"github.com/razshare/frizzante/libcon"
	"github.com/razshare/frizzante/libnotifier"
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
	Notifier        *libnotifier.Notifier
	Efs             embed.FS
	ViewRoot        string
	PublicRoot      string
	ViewServer      string
	ViewIndex       string
	Upgrader        *websocket.Upgrader
	Guards          []Guard
}

type Guard struct {
	Name    string
	Handler func(c *libcon.Connection, allow func())
	Tags    []string
}

type Route struct {
	Pattern string
	Handler func(c *libcon.Connection)
	Tags    []string
}
