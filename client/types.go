package client

import (
	"embed"
	"github.com/gorilla/websocket"
	"github.com/razshare/frizzante/view"
	"log"
	"net/http"
)

type Scope struct {
	Render     view.Render
	ErrorLog   *log.Logger
	InfoLog    *log.Logger
	WebSocket  *websocket.Conn
	Efs        embed.FS
	PublicRoot string
	SessionId  string
	EventName  string
	EventId    int64
	Locked     bool
	Status     int
}

type Client struct {
	Writer  http.ResponseWriter
	Request *http.Request
	Scope   Scope
}
