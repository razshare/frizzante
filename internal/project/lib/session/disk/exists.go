package session

import (
	"main/lib/core/client"
	"main/lib/core/files"
	"main/lib/core/receive"
	"path/filepath"
)

func Exists(c *client.Client) bool {
	id := receive.SessionId(c)
	mtx := Lock(c)
	defer mtx.Unlock()
	return files.IsFile(filepath.Join(".gen", "sessions", id+".json"))
}
