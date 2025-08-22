package send

import (
	"fmt"
	"github.com/razshare/frizzante/client"
	"net/url"
)

// Cookie sends a cookies to the client.
func Cookie(client *client.Client, key string, value string) {
	Header(client, "Set-Cookie", fmt.Sprintf("%s=%s; Path=/; HttpOnly", url.QueryEscape(key), url.QueryEscape(value)))
}
