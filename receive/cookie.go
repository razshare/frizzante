package receive

import (
	"github.com/razshare/frizzante/client"
	"github.com/razshare/frizzante/stack"
	"net/url"
)

// Cookie reads the contents of a cookie from the message and returns the value.
//
// Compatible with web sockets.
func Cookie(c *client.Client, key string) string {
	cookie, cookieError := c.Request.Cookie(key)
	if cookieError != nil {
		c.Scope.ErrorLog.Println(cookieError, stack.Trace())
		return ""
	}

	data, queryError := url.QueryUnescape(cookie.Value)
	if queryError != nil {
		c.Scope.ErrorLog.Println(queryError, stack.Trace())
		return ""
	}

	return data
}
