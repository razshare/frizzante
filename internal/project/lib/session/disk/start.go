package session

import (
	"main/lib/core/client"
	"sync"
)

var Mutexes = map[string]*sync.Mutex{}

func Start(c *client.Client) *State {
	if !Exists(c) {
		s := New()
		Save(c, s)
		return s
	}

	return Load(c)
}
