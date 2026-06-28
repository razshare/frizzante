package scopes

import (
	"embed"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/razshare/frizzante/internal/project/lib/core/databases/schema"
	"github.com/razshare/frizzante/internal/project/lib/core/views/renders"
)

type Http struct {
	Request   http.Request
	Efs       embed.FS
	Render    renders.Render
	ErrorLog  *log.Logger
	InfoLog   *log.Logger
	Writer    http.ResponseWriter
	SessionId string
	EventName string
	EventId   int64
	Status    int
	WebSocket *websocket.Conn
	Locked    bool
	Parsed    bool
	Queries   *schema.Queries
}
