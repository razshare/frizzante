package sessions

import (
	"github.com/razshare/frizzante/archives"
	"github.com/razshare/frizzante/servers"
)

type Session[T any] struct {
	Archive    archives.Archive
	Connection *servers.Connection
	State      *T
}
