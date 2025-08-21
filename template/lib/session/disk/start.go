//gen:mod "disk" "session"
package disk

import (
	"encoding/json"
	"github.com/razshare/frizzante/client"
	"github.com/razshare/frizzante/files"
	"github.com/razshare/frizzante/receive"
	"github.com/razshare/frizzante/stack"
	"os"
	"path/filepath"
	"sync"
)

//gen:mod "mutexes" "Mutexes"
var mutexes = map[string]*sync.Mutex{}

//gen:mod "start" "Start"
//gen:mod "state" "State"
func start(c *client.Client) *state {
	//gen:mod "exists" "Exists"
	if !exists(c) {
		//gen:mod "newState" "New"
		s := newState()
		//gen:mod "save" "Save"
		save(c, s)
		return s
	}

	//gen:mod "load" "Load"
	return load(c)
}

//gen:mod "exists" "Exists"
func exists(c *client.Client) bool {
	id := receive.SessionId(c)
	//gen:mod "lock" "Lock"
	mtx := lock(c)
	defer mtx.Unlock()
	return files.IsFile(filepath.Join(".gen", "sessions", id+".json"))
}

//gen:mod "save" "Save"
//gen:mod "state" "State"
func save(c *client.Client, s *state) {
	//gen:mod "lock" "Lock"
	mtx := lock(c)
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

//gen:mod "load" "Load"
//gen:mod "state" "State"
func load(c *client.Client) *state {
	//gen:mod "lock" "Lock"
	mtx := lock(c)
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

	//gen:mod "newState" "New"
	v := newState()

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

//gen:mod "lock" "Lock"
func lock(c *client.Client) *sync.Mutex {
	id := receive.SessionId(c)
	//gen:mod "mutexes" "Mutexes"
	mtx, ok := mutexes[id]

	if !ok {
		mtx = &sync.Mutex{}
		mtx.Lock()
		//gen:mod "mutexes" "Mutexes"
		mutexes[id] = mtx
	} else {
		mtx.Lock()
	}

	return mtx
}
