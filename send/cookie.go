package send

import (
	"fmt"
	"github.com/razshare/frizzante/client"
	"net/url"
)

// Cookie sends a cookies to the client.
func Cookie(c *client.Client, key string, value string) {
	Header(c, "Set-Cookie", fmt.Sprintf("%s=%s; Path=/; HttpOnly", url.QueryEscape(key), url.QueryEscape(value)))
}
