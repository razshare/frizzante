package send

import (
	"encoding/json"
	"net/http"
)

// Json sends json content.
//
// Compatible with web sockets and server sent events.
func Json(writer http.ResponseWriter, value any) (err error) {
	var data []byte
	if data, err = json.Marshal(value); err != nil {
		return
	}
	_, err = writer.Write(data)
	return
}
