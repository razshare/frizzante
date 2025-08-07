package server

import (
	"embed"
	"github.com/gorilla/websocket"
	"github.com/razshare/frizzante/container"
	"log"
	"net/http"
)

type Server struct {
	http.Server
	Container   container.Configuration
	Guards      []Guard
	Routes      []Route
	InfoLog     *log.Logger
	Efs         embed.FS
	Stop        chan any
	SecureAddr  string
	Certificate string
	PublicRoot  string
	Dotenv      string
	Key         string
}

type Connection struct {
	Status    int
	Locked    bool
	EventId   int64
	EventName string
	Request   *http.Request
	Writer    http.ResponseWriter
	WebSocket *websocket.Conn
	Container *container.Container
	Server    *Server
}

type Route struct {
	Pattern string
	Handler func(c *Connection)
	Tags    []string
}

type Guard struct {
	Name    string
	Handler func(c *Connection, pass func())
	Tags    []string
}
