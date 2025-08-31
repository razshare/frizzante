package session

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"

	"github.com/razshare/frizzante/client"
	"github.com/razshare/frizzante/files"
	"github.com/razshare/frizzante/receive"
	"github.com/razshare/frizzante/stack"
)

var Mutexes = map[string]*sync.Mutex{}

func Start(c *client.Client) *State {
	if !Exists(c) {
		s := New()
		Save(c, s)
		return s
	}

	return Load(c)
}

func Exists(c *client.Client) bool {
	id := receive.SessionId(c)
	mtx := Lock(c)
	defer mtx.Unlock()
	return files.IsFile(filepath.Join(".gen", "sessions", id+".json"))
}

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

func Lock(c *client.Client) *sync.Mutex {
	id := receive.SessionId(c)
	mtx, ok := Mutexes[id]

	if !ok {
		mtx = &sync.Mutex{}
		mtx.Lock()
		Mutexes[id] = mtx
	} else {
		mtx.Lock()
	}

	return mtx
}
