package session

import (
	"sync"

	"github.com/razshare/frizzante/internal/project/lib/core/client"
	"github.com/razshare/frizzante/internal/project/lib/core/receive"
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
