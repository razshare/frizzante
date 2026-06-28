package send

import "github.com/razshare/frizzante/internal/project/lib/core/scopes"

// Navigate redirects the request to a location with status 302.
func Navigate(http *scopes.Http, location string) {
	Redirect(http, location, 302)
	Message(http, "")
}
