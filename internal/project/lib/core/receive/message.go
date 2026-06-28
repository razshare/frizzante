package receive

import (
	"io"

	"github.com/razshare/frizzante/internal/project/lib/core/logs"
	"github.com/razshare/frizzante/internal/project/lib/core/scopes"
	"github.com/razshare/frizzante/internal/project/lib/core/stack"
)

// Message reads the contents of the message and returns the value.
//
// Compatible with web sockets.
func Message(http *scopes.Http) string {
	if http.WebSocket != nil {
		_, data, err := http.WebSocket.ReadMessage()
		if err != nil {
			logs.Errorf(
				http,
				"receive.Message: failed to read WebSocket message: %v\n%s",
				err,
				stack.Trace(),
			)
			return ""
		}
		return string(data)
	}
	data, err := io.ReadAll(http.Request.Body)
	if err != nil {
		logs.Errorf(
			http,
			"receive.Message: failed to read request body: %v\n%s",
			err,
			stack.Trace(),
		)
		return ""
	}
	return string(data)
}
