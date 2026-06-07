package receive

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/razshare/frizzante/internal/project/lib/core/clients"
	"github.com/razshare/frizzante/internal/project/lib/core/files"
	"github.com/razshare/frizzante/internal/project/lib/core/logs"
	"github.com/razshare/frizzante/internal/project/lib/core/stack"
)

func SessionActions(client *clients.Client) (load func(value any) bool, save func(value any) bool) {
	id := SessionId(client)
	baseDirectory := filepath.Join(".gen", "sessions")
	fileName := filepath.Join(baseDirectory, id+".json")
	load = func(value any) bool {
		if files.IsFile(fileName) {
			var err error
			var data []byte
			if data, err = os.ReadFile(fileName); err != nil {
				logs.Errorf(
					client,
					"receive.SessionActions: failed to read session file: %v\n%s",
					err,
					stack.Trace(),
				)
				return false
			}
			if err = json.Unmarshal(data, value); err != nil {
				logs.Errorf(
					client,
					"receive.SessionActions: failed to unmarshal session data: %v\n%s",
					err,
					stack.Trace(),
				)
				return false
			}
		}
		return true
	}
	save = func(value any) bool {
		var err error
		var data []byte
		if !files.IsDirectory(baseDirectory) {
			if err = os.MkdirAll(baseDirectory, os.ModePerm); err != nil {
				logs.Errorf(
					client,
					"receive.SessionActions: failed to create sessions directory: %v\n%s",
					err,
					stack.Trace(),
				)
				return false
			}
		}
		if data, err = json.MarshalIndent(value, "", "    "); err != nil {
			logs.Errorf(
				client,
				"receive.SessionActions: failed to marshal session data: %v\n%s",
				err,
				stack.Trace(),
			)
			return false
		}
		if err = os.WriteFile(fileName, data, os.ModePerm); err != nil {
			logs.Errorf(
				client,
				"receive.SessionActions: failed to write session file: %v\n%s",
				err,
				stack.Trace(),
			)
			return false
		}
		return true
	}
	return
}
