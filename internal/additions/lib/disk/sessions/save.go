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

func Save(client *clients.Client) {
	id := receive.SessionId(client)
	fileName := filepath.Join(DirectoryName, id+".json")

	var exists bool
	var mutex *sync.Mutex
	if mutex, exists = Mutexes[id]; !exists {
		mutex = &sync.Mutex{}
		Mutexes[id] = mutex
	}

	var session *Session
	if session, exists = Sessions[id]; !exists {
		session = New()
	}

	mutex.Lock()
	defer mutex.Unlock()

	if !files.IsDirectory(DirectoryName) {
		if err := os.MkdirAll(DirectoryName, os.ModePerm); err != nil {
			client.Config.ErrorLog.Println(err, stack.Trace())
			return
		}
	}
	if data, err := json.MarshalIndent(session, "", "    "); err != nil {
		client.Config.ErrorLog.Println(err, stack.Trace())
		return
	} else if err = os.WriteFile(fileName, data, os.ModePerm); err != nil {
		client.Config.ErrorLog.Println(err, stack.Trace())
		return
	}
}
