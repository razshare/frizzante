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

var Sessions = map[string]*Session{}
var Mutexes = map[string]*sync.Mutex{}
var DirectoryName = filepath.Join(".gen", "sessions")

func Sync(client *clients.Client) (pull func() (session *Session), push func(session *Session)) {
	id := receive.SessionId(client)
	fileName := filepath.Join(DirectoryName, id+".json")
	if _, exists := Mutexes[id]; !exists {
		Mutexes[id] = &sync.Mutex{}
	}
	pull = func() (session *Session) {
		if mutex, exists := Mutexes[id]; exists {
			mutex.Lock()
			defer mutex.Unlock()
		} else {
			client.Config.ErrorLog.Println("could not save session because relative mutex is missing", stack.Trace())
			return
		}

		if files.IsFile(fileName) {
			session = &Session{}
			if data, err := os.ReadFile(fileName); err != nil {
				client.Config.ErrorLog.Println(err, stack.Trace())
				return
			} else if err = json.Unmarshal(data, session); err != nil {
				client.Config.ErrorLog.Println(err, stack.Trace())
				return
			}
		} else {
			session = New()
		}

		Sessions[id] = session
		return
	}
	push = func(session *Session) {
		if mutex, exists := Mutexes[id]; exists {
			mutex.Lock()
			defer mutex.Unlock()
		} else {
			client.Config.ErrorLog.Println("could not save session because relative mutex is missing", stack.Trace())
			return
		}

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
	return
}
