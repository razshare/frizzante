package servers

import (
	"embed"
	"github.com/gorilla/websocket"
	"log"
	"net"
	"net/http"
)

type Server struct {
	Guards        []Guard
	Routes        []Route
	Efs           embed.FS
	Connections   map[string]*net.Conn
	InfoLog       *log.Logger
	Address       string
	SecureAddress string
	PublicRoot    string
	AppRoot       string
	ServerJs      string
	IndexHtml     string
	Dotenv        string
	Certificate   string
	Key           string
	http.Server
}

type Connection struct {
	WebSocket *websocket.Conn
	Web       *Server
	Request   *http.Request
	Writer    http.ResponseWriter
	Status    int
	EventId   int64
	EventName string
	SessionId string
	Locked    bool
}

type Route struct {
	Pattern string
	Handler func(c *Connection)
	Tags    []string
}

type Guard struct {
	Name    string
	Handler func(c *Connection, allow func())
	Tags    []string
}
