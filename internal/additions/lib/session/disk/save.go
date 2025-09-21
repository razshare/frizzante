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

func Save(client *client.Client, s *State) {
	mtx := Lock(client)
	defer mtx.Unlock()

	dname := filepath.Join(".gen", "sessions")
	if !files.IsDirectory(dname) {
		err := os.MkdirAll(dname, os.ModePerm)
		if err != nil {
			client.Config.ErrorLog.Println(err, stack.Trace())
			return
		}
	}

	id := receive.SessionId(client)

	fname := filepath.Join(dname, id+".json")

	data, err := json.Marshal(s)
	if err != nil {
		client.Config.ErrorLog.Println(err, stack.Trace())
		return
	}

	if err = os.WriteFile(fname, data, os.ModePerm); err != nil {
		client.Config.ErrorLog.Println(err, stack.Trace())
	}
}
