package sessions

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/razshare/frizzante/internal/project/lib/core/clients"
	"github.com/razshare/frizzante/internal/project/lib/core/files"
	"github.com/razshare/frizzante/internal/project/lib/core/receive"
	"github.com/razshare/frizzante/internal/project/lib/core/stacks"
)

func Save(client *clients.Client, s *Session) {
	mtx := Lock(client)
	defer mtx.Unlock()

	dname := filepath.Join(".gen", "sessions")
	if !files.IsDirectory(dname) {
		err := os.MkdirAll(dname, os.ModePerm)
		if err != nil {
			client.Config.ErrorLog.Println(err, stacks.Trace())
			return
		}
	}

	id := receive.SessionId(client)

	fname := filepath.Join(dname, id+".json")

	data, err := json.Marshal(s)
	if err != nil {
		client.Config.ErrorLog.Println(err, stacks.Trace())
		return
	}

	if err = os.WriteFile(fname, data, os.ModePerm); err != nil {
		client.Config.ErrorLog.Println(err, stacks.Trace())
	}
}
