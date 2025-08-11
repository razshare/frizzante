package client

import (
	"github.com/gorilla/websocket"
	"github.com/razshare/frizzante/container"
	"net/http"
)

type Scope struct {
	Container *container.Container
	WebSocket *websocket.Conn
	SessionId string
	EventName string
	EventId   int64
	Locked    bool
	Status    int
}

type Client struct {
	Writer  http.ResponseWriter
	Request *http.Request
	Scope   Scope
}
