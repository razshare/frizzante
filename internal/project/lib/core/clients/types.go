package clients

import (
	"embed"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type Client struct {
	SessionId string
	EventName string
	EventId   int64
	Sink      []func()
	Status    int
	Config    *Config
	Request   *http.Request
	WebSocket *websocket.Conn
	Writer    http.ResponseWriter
	Locked    bool
	Parsed    bool
}

type Config struct {
	Efs      embed.FS
	ErrorLog *log.Logger
	InfoLog  *log.Logger
}
