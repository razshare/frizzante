package client

import (
	"embed"
	"log"
	_view "main/lib/core/view"
	"net/http"

	"github.com/gorilla/websocket"
)

type Client struct {
	Config    *Config
	Request   *http.Request
	WebSocket *websocket.Conn
	Writer    http.ResponseWriter
	SessionId string
	EventName string
	EventId   int64
	Locked    bool
	Status    int
}

type Config struct {
	PublicRoot string
	Efs        embed.FS
	ErrorLog   *log.Logger
	InfoLog    *log.Logger
	Render     func(view _view.View) (html string, err error)
}
