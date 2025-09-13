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
