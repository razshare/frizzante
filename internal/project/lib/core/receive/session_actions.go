package receive

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"

	"github.com/razshare/frizzante/internal/project/lib/core/clients"
	"github.com/razshare/frizzante/internal/project/lib/core/files"
	"github.com/razshare/frizzante/internal/project/lib/core/stack"
)

func SessionActions(client *clients.Client) (load func(value any) bool, save func(value any) bool) {
	id := SessionId(client)
	baseDirectory := filepath.Join(".gen", "sessions")
	fileName := filepath.Join(baseDirectory, id+".json")

	var exists bool
	var mutex *sync.Mutex
	if mutex, exists = Mutexes[id]; !exists {
		mutex = &sync.Mutex{}
		Mutexes[id] = mutex
	}

	load = func(value any) bool {
		if files.IsFile(fileName) {
			var err error
			var data []byte
			if data, err = os.ReadFile(fileName); err != nil {
				client.Options.ErrorLog.Println(err, stack.Trace())
				return false
			}
			if err = json.Unmarshal(data, value); err != nil {
				client.Options.ErrorLog.Println(err, stack.Trace())
				return false
			}
		}
		return true
	}

	save = func(value any) bool {
		mutex.Lock()
		defer mutex.Unlock()
		var err error
		var data []byte

		if !files.IsDirectory(baseDirectory) {
			if err = os.MkdirAll(baseDirectory, os.ModePerm); err != nil {
				client.Options.ErrorLog.Println(err, stack.Trace())
				return false
			}
		}

		if data, err = json.MarshalIndent(value, "", "    "); err != nil {
			client.Options.ErrorLog.Println(err, stack.Trace())
			return false
		}

		if err = os.WriteFile(fileName, data, os.ModePerm); err != nil {
			client.Options.ErrorLog.Println(err, stack.Trace())
			return false
		}
		return true
	}
	return
}
