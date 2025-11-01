package sessions

import (
	"github.com/razshare/frizzante/internal/project/lib/core/clients"
	"github.com/razshare/frizzante/internal/project/lib/core/receive"
)

var Sessions = map[string]*Session{}

func Start(client *clients.Client) (session *Session) {
	id := receive.SessionId(client)
	var exists bool
	if session, exists = Sessions[id]; !exists {
		session = New()
		Sessions[id] = session
		return
	}
	return
}
