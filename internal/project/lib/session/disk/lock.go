package session

import (
	"main/lib/core/client"
	"main/lib/core/receive"
	"sync"
)

func Lock(c *client.Client) *sync.Mutex {
	id := receive.SessionId(c)
	mtx, ok := Mutexes[id]

	if !ok {
		mtx = &sync.Mutex{}
		mtx.Lock()
		Mutexes[id] = mtx
	} else {
		mtx.Lock()
	}

	return mtx
}
