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
