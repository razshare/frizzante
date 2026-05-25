package receive

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/razshare/frizzante/internal/project/lib/core/clients"
	"github.com/razshare/frizzante/internal/project/lib/core/files"
	"github.com/razshare/frizzante/internal/project/lib/core/stack"
)

func Session(client *clients.Client, value any) bool {
	id := SessionId(client)
	baseDirectory := filepath.Join(".gen", "sessions")
	fileName := filepath.Join(baseDirectory, id+".json")
	if files.IsFile(fileName) {
		var err error
		var data []byte
		if data, err = os.ReadFile(fileName); err != nil {
			client.Options.ErrorLog.Printf(
				"receive.Session: failed to read session file: %v\n%s",
				err,
				stack.Trace(),
			)
			return false
		}
		if err = json.Unmarshal(data, value); err != nil {
			client.Options.ErrorLog.Printf(
				"receive.Session: failed to unmarshal session data: %v\n%s",
				err,
				stack.Trace(),
			)
			return false
		}
	}
	go func() {
		<-client.Request.Context().Done()
		var err error
		var data []byte
		if !files.IsDirectory(baseDirectory) {
			if err = os.MkdirAll(baseDirectory, os.ModePerm); err != nil {
				client.Options.ErrorLog.Printf(
					"receive.Session: failed to create sessions directory: %v\n%s",
					err,
					stack.Trace(),
				)
				return
			}
		}
		if data, err = json.MarshalIndent(value, "", "    "); err != nil {
			client.Options.ErrorLog.Printf(
				"receive.Session: failed to marshal session data: %v\n%s",
				err,
				stack.Trace(),
			)
			return
		}
		if err = os.WriteFile(fileName, data, os.ModePerm); err != nil {
			client.Options.ErrorLog.Printf(
				"receive.Session: failed to write session file: %v\n%s",
				err,
				stack.Trace(),
			)
			return
		}
	}()
	return true
}
