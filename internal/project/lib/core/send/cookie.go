package send

import (
	"fmt"
	"net/url"

	_client "github.com/razshare/frizzante/internal/project/lib/core/client"
)

// Cookie sends a cookies to the client.
func Cookie(client *_client.Client, key string, value string) {
	Header(client, "Set-Cookie", fmt.Sprintf("%s=%s; Path=/; HttpOnly", url.QueryEscape(key), url.QueryEscape(value)))
}
