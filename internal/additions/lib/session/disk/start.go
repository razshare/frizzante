package session

import (
	"sync"

	"github.com/razshare/frizzante/internal/project/lib/core/client"
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
