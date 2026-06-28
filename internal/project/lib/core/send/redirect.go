package send

import (
	"errors"
	"net/http"
)

// Redirect redirects the request to a location with a status.
func Redirect(writer http.ResponseWriter, location string, status int) (err error) {
	header := writer.Header()
	header.Set("Location", location)
	writer.WriteHeader(status)
	var ok bool
	var flusher http.Flusher
	if flusher, ok = writer.(http.Flusher); !ok {
		err = errors.New("send.Redirect: could not retrieve flusher")
		return
	}
	flusher.Flush()
	return
}
