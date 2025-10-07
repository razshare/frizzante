package session

import (
	"path/filepath"

	_client "github.com/razshare/frizzante/internal/project/lib/core/client"
	"github.com/razshare/frizzante/internal/project/lib/core/files"
	"github.com/razshare/frizzante/internal/project/lib/core/receive"
)

func Exists(client *_client.Client) bool {
	id := receive.SessionId(client)
	mtx := Lock(client)
	defer mtx.Unlock()
	return files.IsFile(filepath.Join(".gen", "sessions", id+".json"))
}
