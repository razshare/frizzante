package sessions

import (
	"sync"

	"github.com/razshare/frizzante/internal/project/lib/core/clients"
)

var Mutexes = map[string]*sync.Mutex{}

func Start(client *clients.Client) *Session {
	if !Exists(client) {
		state := New()
		Save(client, state)
		return state
	}

	return Load(client)
}
