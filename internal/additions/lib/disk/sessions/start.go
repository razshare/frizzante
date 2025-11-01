package sessions

import "github.com/razshare/frizzante/internal/project/lib/core/clients"

func Start(client *clients.Client) (session *Session) {
	pull, push := Sync(client)
	client.Deferred = append(client.Deferred, func() { push(session) })
	session = pull()
	return
}
