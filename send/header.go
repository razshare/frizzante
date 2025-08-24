package send

import (
	"github.com/razshare/frizzante/client"
	"github.com/razshare/frizzante/stack"
)

// Header sends a header field.
//
// If the status has not been sent already, a default "200 OK" status will be sent immediately.
//
// This means the status will become locked and further attempts to send the status will fail with an error.
//
// All errors are sent to the server notifier.
func Header(c *client.Client, k string, v string) {
	if c.Locked {
		c.Config.ErrorLog.Println("header is locked", stack.Trace())
		return
	}

	c.Writer.Header().Set(k, v)
}

// Headers sends header fields.
func Headers(c *client.Client, f map[string]string) {
	if c.Locked {
		c.Config.ErrorLog.Println("header is locked", stack.Trace())
		return
	}

	for k, v := range f {
		c.Writer.Header().Set(k, v)
	}
}

// Redirect redirects the request to a location with a status.
func Redirect(c *client.Client, l string, s int) {
	Status(c, s)
	Header(c, "Location", l)
}

// Navigate redirects the request to a location with status 302.
func Navigate(c *client.Client, l string) {
	Redirect(c, l, 302)
	Message(c, "")
}

// ContentType sets the Content-Type header field.
func ContentType(c *client.Client, t string) {
	Header(c, "Content-Type", t)
}
