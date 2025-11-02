package receive

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/razshare/frizzante/internal/project/lib/core/clients"
	"github.com/razshare/frizzante/internal/project/lib/core/files"
	"github.com/razshare/frizzante/internal/project/lib/core/stack"
)

func SessionSaver(client *clients.Client) (save func(value any) bool) {
	baseDirectory := filepath.Join(".gen", "sessions")
	fileName := filepath.Join(baseDirectory, SessionId(client)+".json")
	save = func(value any) bool {
		var err error
		if !files.IsDirectory(baseDirectory) {
			if err = os.MkdirAll(baseDirectory, os.ModePerm); err != nil {
				client.Config.ErrorLog.Println(err, stack.Trace())
				return false
			}
		}
		var data []byte
		if data, err = json.MarshalIndent(value, "", "    "); err != nil {
			client.Config.ErrorLog.Println(err, stack.Trace())
			return false
		}
		if err = os.WriteFile(fileName, data, os.ModePerm); err != nil {
			client.Config.ErrorLog.Println(err, stack.Trace())
			return false
		}
		return true
	}
	return
}
