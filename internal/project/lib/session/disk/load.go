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

func Load(c *client.Client) *State {
	mtx := Lock(c)
	defer mtx.Unlock()

	dn := filepath.Join(".gen", "sessions")
	if !files.IsDirectory(dn) {
		err := os.MkdirAll(dn, os.ModePerm)
		if err != nil {
			c.Config.ErrorLog.Println(err, stack.Trace())
			return nil
		}
	}

	id := receive.SessionId(c)
	n := filepath.Join(dn, id+".json")

	v := New()

	var d []byte
	d, err := os.ReadFile(n)
	if err != nil {
		c.Config.ErrorLog.Println(err, stack.Trace())
		return v
	}

	err = json.Unmarshal(d, v)
	if err != nil {
		c.Config.ErrorLog.Println(err, stack.Trace())
		return v
	}
	return v
}
