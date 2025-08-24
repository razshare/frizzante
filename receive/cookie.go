package receive

import (
	"github.com/razshare/frizzante/client"
	"github.com/razshare/frizzante/stack"
	"net/url"
)

// Cookie reads the contents of a cookie from the message and returns the value.
//
// Compatible with web sockets.
func Cookie(c *client.Client, k string) string {
	ck, err := c.Request.Cookie(k)
	if err != nil {
		c.Config.ErrorLog.Println(err, stack.Trace())
		return ""
	}

	d, err := url.QueryUnescape(ck.Value)
	if err != nil {
		c.Config.ErrorLog.Println(err, stack.Trace())
		return ""
	}

	return d
}
