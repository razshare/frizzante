package receive

import (
	"github.com/razshare/frizzante/client"
)

// Header reads a header field and returns the value.
//
// Compatible with web sockets.
func Header(c *client.Client, k string) string {
	return c.Request.Header.Get(k)
}

// ContentType reads the Content-Type header field and returns the value.
//
// Compatible with web sockets.
func ContentType(c *client.Client) string {
	return c.Request.Header.Get("Content-Type")
}

// Accept reads if the Accept header entries and returns the values.
func Accept(c *client.Client) string {
	return c.Request.Header.Get("Accept")
}
