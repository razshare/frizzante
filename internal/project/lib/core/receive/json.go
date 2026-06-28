package receive

import (
	"encoding/json"
	"io"

	"github.com/razshare/frizzante/internal/project/lib/core/logs"
	"github.com/razshare/frizzante/internal/project/lib/core/scopes"
	"github.com/razshare/frizzante/internal/project/lib/core/stack"
)

// Json reads the next JSON-encoded message from the
// client and stores it in the value pointed to by value.
//
// Compatible with web sockets and server sent events.
func Json(http *scopes.Http, value any) bool {
	if http.WebSocket != nil {
		if err := http.WebSocket.ReadJSON(&value); err != nil {
			logs.Errorf(
				http,
				"receive.Json: failed to read WebSocket JSON message: %v\n%s",
				err,
				stack.Trace(),
			)
			return false
		}
		return true
	}
	data, err := io.ReadAll(http.Request.Body)
	if err != nil {
		logs.Errorf(
			http,
			"receive.Json: failed to read request body: %v\n%s",
			err,
			stack.Trace(),
		)
		return false
	}
	if err = json.Unmarshal(data, &value); err != nil {
		logs.Errorf(
			http,
			"receive.Json: failed to unmarshal JSON: %v\n%s",
			err,
			stack.Trace(),
		)
	}
	return true
}
