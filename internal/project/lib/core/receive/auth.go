package receive

import "github.com/razshare/frizzante/internal/project/lib/core/client"

// BasicAuth reads the username and password provided
// in the request's Authorization header and stores them into the value
// pointed to by username and password, if the request uses HTTP Basic Authentication.
//
// See RFC 2617, Section 2
func BasicAuth(client *client.Client) (username string, password string) {
	username, password, _ = client.Request.BasicAuth()
	return
}
