package receive

import "github.com/razshare/frizzante/client"

// BasicAuth returns the username and password provided
// in the request's Authorization header, if the request
// uses HTTP Basic Authentication. See RFC 2617, Section 2
func BasicAuth(c *client.Client) (username string, password string, ok bool) {
	return c.Request.BasicAuth()
}
