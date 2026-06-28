package receive

import (
	"net/http"
	"net/url"
)

// Cookie reads the contents of a cookie from the message and returns the value.
func Cookie(request *http.Request, key string) (value string, err error) {
	var cookie *http.Cookie
	if cookie, err = request.Cookie(key); err != nil {
		return
	}
	value, err = url.QueryUnescape(cookie.Value)
	return
}
