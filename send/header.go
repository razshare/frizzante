package send

import (
	"fmt"
	"github.com/razshare/frizzante/conn"
	"github.com/razshare/frizzante/stack"
	"net/url"
)

// Navigate redirects the request to a location with status 302.
func Navigate(c *conn.Conn, l string) {
	Redirect(c, l, 302)
	Flush(c)
}

// Redirect redirects the request to a location with a status.
func Redirect(c *conn.Conn, l string, status int) {
	Status(c, status)
	Header(c, "Location", l)
}

// Status sets the status code.
//
// This will lock the status, which makes it
// so that the increaseIndex time you invoke this
// function it will fail with an error.
//
// All errors are sent to the server notifier.
func Status(c *conn.Conn, s int) {
	if c.Locked {
		c.Container.Config.ErrorLog.Println("status is locked", stack.Trace())
		return
	}

	c.Status = s
}

// Header sends a header field.
//
// If the status has not been sent already, a default "200 OK" status will be sent immediately.
//
// This means the status will become locked and further attempts to send the status will fail with an error.
//
// All errors are sent to the server notifier.
func Header(c *conn.Conn, k string, v string) {
	if c.Locked {
		c.Container.Config.ErrorLog.Println("header is locked", stack.Trace())
		return
	}

	c.Writer.Header().Set(k, v)
}

// Headers sends header fields.
func Headers(c *conn.Conn, h map[string]string) {
	if c.Locked {
		c.Container.Config.ErrorLog.Println("header is locked", stack.Trace())
		return
	}

	for key, value := range h {
		c.Writer.Header().Set(key, value)
	}
}

// ContentType sets the Content-Type header field.
func ContentType(c *conn.Conn, t string) {
	Header(c, "Content-Type", t)
}

// Cookie sends a cookies to the client.
func Cookie(c *conn.Conn, key string, value string) {
	Header(c, "Set-Cookie", fmt.Sprintf("%s=%s; Path=/; HttpOnly", url.QueryEscape(key), url.QueryEscape(value)))
}
