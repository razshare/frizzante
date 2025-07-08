package sessions

import (
	"github.com/razshare/frizzante/archives"
	"github.com/razshare/frizzante/connections"
)

type SessionStarter[T any] = func(con *connections.Connection, state T) *Session[T]

type Session[T any] struct {
	Archive    archives.Archive
	Connection *connections.Connection
	State      *T
}
