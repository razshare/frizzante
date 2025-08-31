package send

import (
	"encoding/json"
	"main/lib/core/client"

	"github.com/razshare/frizzante/stack"
)

// Json sends json content.
//
// If the status code or the header have not been sent already, a default status of "200 OK" will be sent immediately along with whatever headers you've previously defined.
//
// The status code and the header will become locked and further attempts to send either of them will fail with an error.
//
// All errors are sent to the server notifier.
//
// Compatible with web sockets.
func Json(client *client.Client, value any) {
	data, err := json.Marshal(value)
	if err != nil {
		client.Config.ErrorLog.Println(err, stack.Trace())
		return
	}

	if nil == client.WebSocket {
		if client.Writer.Header().Get("Content-Type") == "" {
			client.Writer.Header().Set("Content-Type", "application/json")
		}
	}

	Content(client, data)
}
