package client

import (
	"embed"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/razshare/frizzante/internal/project/lib/core/view"
)

type Client struct {
	SessionId string
	EventName string
	EventId   int64
	Status    int
	Token     string
	UserId    string
	Config    *Config
	Request   *http.Request
	WebSocket *websocket.Conn
	Writer    http.ResponseWriter
	Locked    bool
	Parsed    bool
}

type Config struct {
	PublicRoot string
	Efs        embed.FS
	ErrorLog   *log.Logger
	InfoLog    *log.Logger
	Render     func(view view.View) (html string, err error)
}
