package session

import (
	"sync"

	_client "github.com/razshare/frizzante/internal/project/lib/core/client"
)

var Mutexes = map[string]*sync.Mutex{}

func Start(client *_client.Client) *State {
	if !Exists(client) {
		state := New()
		Save(client, state)
		return state
	}

	return Load(client)
}
