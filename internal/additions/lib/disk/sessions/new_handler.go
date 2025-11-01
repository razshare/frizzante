package sessions

import "github.com/razshare/frizzante/internal/project/lib/core/clients"

func NewHandler(handle func(client *clients.Client)) func(client *clients.Client) {
	return func(client *clients.Client) {
		session := Start(client)
		defer Save(session, client)
		handle(client)
	}
}
