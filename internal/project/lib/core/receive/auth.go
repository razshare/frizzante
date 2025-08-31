package receive

import "main/lib/core/client"

// BasicAuth returns the username and password provided
// in the request's Authorization header, if the request
// uses HTTP Basic Authentication. See RFC 2617, Section 2
func BasicAuth(client *client.Client) (user string, pass string, ok bool) {
	return client.Request.BasicAuth()
}
