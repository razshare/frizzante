package sessions

import (
	"github.com/razshare/frizzante/connections"
)

type Session[T any] struct {
	Connection *connections.Connection
	State      *T
}
