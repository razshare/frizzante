package session

import (
	"encoding/json"
	"main/lib/core/client"
	"main/lib/core/files"
	"main/lib/core/receive"
	"main/lib/core/stack"
	"os"
	"path/filepath"
)

func Save(c *client.Client, s *State) {
	mtx := Lock(c)
	defer mtx.Unlock()

	dn := filepath.Join(".gen", "sessions")
	if !files.IsDirectory(dn) {
		err := os.MkdirAll(dn, os.ModePerm)
		if err != nil {
			c.Config.ErrorLog.Println(err, stack.Trace())
			return
		}
	}

	id := receive.SessionId(c)

	n := filepath.Join(dn, id+".json")

	d, err := json.Marshal(s)
	if err != nil {
		c.Config.ErrorLog.Println(err, stack.Trace())
		return
	}

	err = os.WriteFile(n, d, os.ModePerm)
	if err != nil {
		c.Config.ErrorLog.Println(err, stack.Trace())
	}
}
