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
	baseDirectory := filepath.Join(".gen", "sessions")
	var err error
	if !files.IsDirectory(baseDirectory) {
		if err = os.MkdirAll(baseDirectory, os.ModePerm); err != nil {
			client.Config.ErrorLog.Println(err, stack.Trace())
			return false
		}
	}

	fileName := filepath.Join(baseDirectory, SessionId(client)+".json")
	var data []byte
	if files.IsFile(fileName) {
		if data, err = os.ReadFile(fileName); err != nil {
			client.Config.ErrorLog.Println(err, stack.Trace())
			return false
		}
		if err = json.Unmarshal(data, value); err != nil {
			client.Config.ErrorLog.Println(err, stack.Trace())
			return false
		}
	} else {
		if data, err = json.MarshalIndent(value, "", "    "); err != nil {
			client.Config.ErrorLog.Println(err, stack.Trace())
			return false
		}
		if err = os.WriteFile(fileName, data, os.ModePerm); err != nil {
			client.Config.ErrorLog.Println(err, stack.Trace())
			return false
		}
	}

	client.Sink = append(client.Sink, func() {
		if data, err = json.MarshalIndent(value, "", "    "); err != nil {
			client.Config.ErrorLog.Println(err, stack.Trace())
			return
		}
		if err = os.WriteFile(fileName, data, os.ModePerm); err != nil {
			client.Config.ErrorLog.Println(err, stack.Trace())
			return
		}
	})

	return true
}
