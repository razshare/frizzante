package libsession

import (
	"github.com/razshare/frizzante/libarchive"
	"github.com/razshare/frizzante/libcon"
)

type SessionOperator interface {
	Id() string
	Exists() bool
	Load(state any)
	Save(state any)
	Destroy()
	WithArchive(archive libarchive.Archive) SessionOperator
}

type SessionHandler[T any] = func(state *T)

type ConnectedSessionOperator struct {
	Archive    libarchive.Archive
	Connection *libcon.Connection
}
