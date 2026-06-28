package send

import (
	"net/http"
)

// Navigate redirects the request to a location with status 302.
func Navigate(writer http.ResponseWriter, location string) error {
	return Redirect(writer, location, 302)
}
