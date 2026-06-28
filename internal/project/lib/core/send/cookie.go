package send

import (
	"fmt"
	"net/url"

	"github.com/razshare/frizzante/internal/project/lib/core/scopes"
)

// Cookie sends a cookies to the http.
func Cookie(http *scopes.Http, key string, value string) {
	Header(http, "Set-Cookie", fmt.Sprintf("%s=%s; Path=/; HttpOnly", url.QueryEscape(key), url.QueryEscape(value)))
}
