package send

import (
	"github.com/razshare/frizzante/client"
	"github.com/razshare/frizzante/stack"
)

// Navigate redirects the request to a location with status 302.
func Navigate(c *client.Client, l string) {
	Redirect(c, l, 302)
	Message(c, "")
}

// Redirect redirects the request to a location with a status.
func Redirect(c *client.Client, l string, status int) {
	Status(c, status)
	Header(c, "Location", l)
}

// Header sends a header field.
//
// If the status has not been sent already, a default "200 OK" status will be sent immediately.
//
// This means the status will become locked and further attempts to send the status will fail with an error.
//
// All errors are sent to the server notifier.
func Header(c *client.Client, k string, v string) {
	if c.Scope.Locked {
		c.Scope.Container.Config.ErrorLog.Println("header is locked", stack.Trace())
		return
	}

	c.Writer.Header().Set(k, v)
}

// Headers sends header fields.
func Headers(c *client.Client, h map[string]string) {
	if c.Scope.Locked {
		c.Scope.Container.Config.ErrorLog.Println("header is locked", stack.Trace())
		return
	}

	for key, value := range h {
		c.Writer.Header().Set(key, value)
	}
}

// ContentType sets the Content-Type header field.
func ContentType(c *client.Client, t string) {
	Header(c, "Content-Type", t)
}
