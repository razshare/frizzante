package send

import "github.com/razshare/frizzante/internal/project/lib/core/scopes"

// Redirect redirects the request to a location with a status.
func Redirect(http *scopes.Http, location string, status int) {
	Status(http, status)
	Header(http, "Location", location)
}
