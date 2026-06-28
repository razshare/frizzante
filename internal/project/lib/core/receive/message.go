package receive

import (
	"io"
	"net/http"
)

// Message reads the contents of the message and returns the value.
//
// Compatible with web sockets.
func Message(request *http.Request) (value string, err error) {
	var data []byte
	if data, err = io.ReadAll(request.Body); err != nil {
		return
	}
	value = string(data)
	return
}
