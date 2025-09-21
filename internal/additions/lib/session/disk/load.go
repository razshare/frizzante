package session

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/razshare/frizzante/internal/project/lib/core/client"
	"github.com/razshare/frizzante/internal/project/lib/core/files"
	"github.com/razshare/frizzante/internal/project/lib/core/receive"
	"github.com/razshare/frizzante/internal/project/lib/core/stack"
)

func Load(client *client.Client) *State {
	mtx := Lock(client)
	defer mtx.Unlock()

	dname := filepath.Join(".gen", "sessions")
	if !files.IsDirectory(dname) {
		err := os.MkdirAll(dname, os.ModePerm)
		if err != nil {
			client.Config.ErrorLog.Println(err, stack.Trace())
			return nil
		}
	}

	id := receive.SessionId(client)
	fname := filepath.Join(dname, id+".json")

	state := New()

	var data []byte
	data, err := os.ReadFile(fname)
	if err != nil {
		client.Config.ErrorLog.Println(err, stack.Trace())
		return state
	}

	err = json.Unmarshal(data, state)
	if err != nil {
		client.Config.ErrorLog.Println(err, stack.Trace())
		return state
	}
	return state
}
