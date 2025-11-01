package sessions

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"

	"github.com/razshare/frizzante/internal/project/lib/core/clients"
	"github.com/razshare/frizzante/internal/project/lib/core/files"
	"github.com/razshare/frizzante/internal/project/lib/core/receive"
	"github.com/razshare/frizzante/internal/project/lib/core/stack"
)

func Start(client *clients.Client) (session *Session) {
	id := receive.SessionId(client)
	fileName := filepath.Join(DirectoryName, id+".json")

	var mutex *sync.Mutex
	var exists bool
	if mutex, exists = Mutexes[id]; !exists {
		mutex = &sync.Mutex{}
		Mutexes[id] = mutex
	}

	mutex.Lock()
	defer mutex.Unlock()

	if files.IsFile(fileName) {
		session = &Session{}
		if data, err := os.ReadFile(fileName); err != nil {
			client.Config.ErrorLog.Println(err, stack.Trace())
			session = New()
			return
		} else if err = json.Unmarshal(data, session); err != nil {
			client.Config.ErrorLog.Println(err, stack.Trace())
			session = New()
			return
		}
	} else {
		session = New()
	}

	Sessions[id] = session

	return
}
