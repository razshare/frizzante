package web

import (
	"embed"
	"github.com/gorilla/websocket"
	"github.com/razshare/frizzante/connections"
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
	ViewRoot        string
	PublicRoot      string
	ViewServer      string
	ViewIndex       string
	WsUpgrader      *websocket.Upgrader
	Guards          []Guard
}

type Route struct {
	Pattern string
	Handler func(c *connections.Connection)
	Tags    []string
}

type Guard struct {
	Name    string
	Handler func(c *connections.Connection, allow func())
	Tags    []string
}
