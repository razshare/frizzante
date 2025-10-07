package session

import (
	"sync"

	_client "github.com/razshare/frizzante/internal/project/lib/core/client"
	"github.com/razshare/frizzante/internal/project/lib/core/receive"
)

func Lock(client *_client.Client) *sync.Mutex {
	id := receive.SessionId(client)
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
