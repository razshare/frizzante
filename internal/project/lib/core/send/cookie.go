package send

import (
	"fmt"
	"net/url"

	"github.com/razshare/frizzante/internal/project/lib/core/clients"
)

// Cookie sends a cookies to the client.
func Cookie(client *clients.Client, key string, value string) {
	Header(client, "Set-Cookie", fmt.Sprintf("%s=%s; Path=/; HttpOnly", url.QueryEscape(key), url.QueryEscape(value)))
}
