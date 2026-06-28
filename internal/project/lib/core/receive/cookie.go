package receive

import (
	"net/url"

	"github.com/razshare/frizzante/internal/project/lib/core/logs"
	"github.com/razshare/frizzante/internal/project/lib/core/scopes"
	"github.com/razshare/frizzante/internal/project/lib/core/stack"
)

// Cookie reads the contents of a cookie from the message and returns the value.
func Cookie(http *scopes.Http, key string) string {
	cookie, err := http.Request.Cookie(key)
	if err != nil {
		logs.Errorf(
			http,
			"receive.Cookie: failed to read cookie %q: %v\n%s",
			key,
			err,
			stack.Trace(),
		)
		return ""
	}

	data, err := url.QueryUnescape(cookie.Value)
	if err != nil {
		logs.Errorf(
			http,
			"receive.Cookie: failed to unescape cookie %q value: %v\n%s",
			key,
			err,
			stack.Trace(),
		)
		return ""
	}

	return data
}
