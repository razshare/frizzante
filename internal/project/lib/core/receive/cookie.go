package receive

import (
	"net/url"

	_client "github.com/razshare/frizzante/internal/project/lib/core/client"
	"github.com/razshare/frizzante/internal/project/lib/core/stack"
)

// Cookie reads the contents of a cookie from the message and returns the value.
func Cookie(client *_client.Client, key string) string {
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
