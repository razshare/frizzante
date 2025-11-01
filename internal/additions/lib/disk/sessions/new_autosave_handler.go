package sessions

import "github.com/razshare/frizzante/internal/project/lib/core/clients"

func NewAutosaveHandler(handle func(client *clients.Client)) func(client *clients.Client) {
	return func(client *clients.Client) {
		defer Save(client)
		handle(client)
	}
}
