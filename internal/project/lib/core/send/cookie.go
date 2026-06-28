package send

import (
	"fmt"
	"net/http"
	"net/url"
)

// Cookie sends a cookies to the http.
func Cookie(writer http.ResponseWriter, key string, value string) {
	header := writer.Header()
	escapedValue := fmt.Sprintf("%s=%s; Path=/; HttpOnly", url.QueryEscape(key), url.QueryEscape(value))
	header.Set("Set-Cookie", escapedValue)
}
