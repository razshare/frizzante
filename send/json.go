package send

import (
	"encoding/json"
	"github.com/razshare/frizzante/client"
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
func Json(c *client.Client, v any) {
	data, jsonError := json.Marshal(v)
	if jsonError != nil {
		c.Config.ErrorLog.Println(jsonError, stack.Trace())
		return
	}

	if nil == c.WebSocket {
		contentType := c.Writer.Header().Get("Content-Type")
		if "" == contentType {
			c.Writer.Header().Set("Content-Type", "application/json")
		}
	}

	Content(c, data)
}
