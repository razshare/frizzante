package send

import (
	"encoding/json"

	"github.com/razshare/frizzante/internal/project/lib/core/clients"
	"github.com/razshare/frizzante/internal/project/lib/core/stacks"
)

// Json sends json content.
//
// Compatible with web sockets and server sent events.
func Json(client *clients.Client, value any) {
	data, err := json.Marshal(value)
	if err != nil {
		client.Config.ErrorLog.Println(err, stacks.Trace())
		return
	}

	if client.WebSocket == nil {
		if client.Writer.Header().Get("Content-Type") == "" {
			client.Writer.Header().Set("Content-Type", "application/json")
		}
	}

	Content(client, data)
}
