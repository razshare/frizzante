package send

import (
	"fmt"
	"github.com/razshare/frizzante/client"
	"net/url"
)

// Cookie sends a cookies to the client.
func Cookie(c *client.Client, k string, v string) {
	Header(c, "Set-Cookie", fmt.Sprintf("%s=%s; Path=/; HttpOnly", url.QueryEscape(k), url.QueryEscape(v)))
}
