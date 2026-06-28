package send

import (
	"encoding/json"

	"github.com/razshare/frizzante/internal/project/lib/core/logs"
	"github.com/razshare/frizzante/internal/project/lib/core/scopes"
	"github.com/razshare/frizzante/internal/project/lib/core/stack"
)

// Json sends json content.
//
// Compatible with web sockets and server sent events.
func Json(http *scopes.Http, value any) {
	data, err := json.Marshal(value)
	if err != nil {
		logs.Errorf(
			http,
			"send.Json: failed to marshal value (type=%T) to JSON: %v\n%s",
			value,
			err,
			stack.Trace(),
		)
		return
	}
	if http.WebSocket == nil {
		if http.Writer.Header().Get("Content-Type") == "" {
			http.Writer.Header().Set("Content-Type", "application/json")
		}
	}
	Content(http, data)
}
