package sessions

import (
	"path/filepath"

	"github.com/razshare/frizzante/internal/project/lib/core/clients"
	"github.com/razshare/frizzante/internal/project/lib/core/files"
	"github.com/razshare/frizzante/internal/project/lib/core/receive"
)

func Exists(client *clients.Client) bool {
	id := receive.SessionId(client)
	mtx := Lock(client)
	defer mtx.Unlock()
	return files.IsFile(filepath.Join(".gen", "sessions", id+".json"))
}
