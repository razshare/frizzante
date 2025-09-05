package receive

import (
	"main/lib/core/client"
	"net/url"

	"main/lib/core/stack"
)

// Cookie reads the contents of a cookie from the message and returns the value.
//
// Compatible with web sockets.
func Cookie(client *client.Client, key string) string {
	cookie, err := client.Request.Cookie(key)
	if err != nil {
		client.Config.ErrorLog.Println(err, stack.Trace())
		return ""
	}

	data, err := url.QueryUnescape(cookie.Value)
	if err != nil {
		client.Config.ErrorLog.Println(err, stack.Trace())
		return ""
	}

	return data
}
