package receive

import "github.com/razshare/frizzante/internal/project/lib/core/scopes"

// BasicAuth reads the username and password provided
// in the request's Authorization header and stores them into the value
// pointed to by username and password, if the request uses HTTP Basic Authentication.
//
// See RFC 2617, Section 2.
func BasicAuth(http *scopes.Http) (username string, password string) {
	username, password, _ = http.Request.BasicAuth()
	return
}
