package send

import (
	"net/http"
)

// Message sends utf-8 safe content.
//
// Compatible with web sockets and server sent events.
func Message(writer http.ResponseWriter, message string) (err error) {
	_, err = writer.Write([]byte(message))
	return
}
