package receive

import (
	"encoding/json"
	"io"
	"net/http"
)

// Json reads the next JSON-encoded message from the
// client and stores it in the value pointed to by value.
//
// Compatible with web sockets and server sent events.
func Json(request *http.Request, value any) (err error) {
	var data []byte
	if data, err = io.ReadAll(request.Body); err != nil {
		return
	}
	if err = json.Unmarshal(data, &value); err != nil {
		return
	}
	return
}
