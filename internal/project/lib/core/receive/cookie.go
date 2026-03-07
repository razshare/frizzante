package receive

import (
	"net/url"

	"github.com/razshare/frizzante/internal/project/lib/core/clients"
	"github.com/razshare/frizzante/internal/project/lib/core/stack"
)

// Cookie reads the contents of a cookie from the message and returns the value.
func Cookie(client *clients.Client, key string) string {
	cookie, err := client.Request.Cookie(key)
	if err != nil {
		client.Options.ErrorLog.Printf(
			"receive.Cookie: failed to read cookie %q: %v\n%s",
			key,
			err,
			stack.Trace(),
		)
		return ""
	}

	data, err := url.QueryUnescape(cookie.Value)
	if err != nil {
		client.Options.ErrorLog.Printf(
			"receive.Cookie: failed to unescape cookie %q value: %v\n%s",
			key,
			err,
			stack.Trace(),
		)
		return ""
	}

	return data
}
