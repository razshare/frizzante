package receive

import (
	"github.com/razshare/frizzante/conn"
	"github.com/razshare/frizzante/stack"
	"net/url"
)

// Path reads a parameters fields and returns the value.
//
// Compatible with web sockets.
func Path(c *conn.Conn, k string) string {
	return c.Request.PathValue(k)
}

// Header reads a header field and returns the value.
//
// Compatible with web sockets.
func Header(c *conn.Conn, k string) string {
	return c.Request.Header.Get(k)
}

// ContentType reads the Content-Type header field and returns the value.
//
// Compatible with web sockets.
func ContentType(c *conn.Conn) string {
	return c.Request.Header.Get("Content-Type")
}

// Accept reads if the Accept header entries and returns the values.
func Accept(c *conn.Conn) string {
	return c.Request.Header.Get("Accept")
}

// Cookie reads the contents of a cookie from the message and returns the value.
//
// Compatible with web sockets.
func Cookie(c *conn.Conn, key string) string {
	cookie, cookieError := c.Request.Cookie(key)
	if cookieError != nil {
		c.Container.Config.ErrorLog.Println(cookieError, stack.Trace())
		return ""
	}

	data, queryError := url.QueryUnescape(cookie.Value)
	if queryError != nil {
		c.Container.Config.ErrorLog.Println(queryError, stack.Trace())
		return ""
	}

	return data
}
