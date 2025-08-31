package client

import (
	"embed"
	"log"
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
	ErrorLog   *log.Logger
	InfoLog    *log.Logger
	PublicRoot string
	Efs        embed.FS
}
