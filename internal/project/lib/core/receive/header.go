package receive

import "github.com/razshare/frizzante/internal/project/lib/core/clients"

// Header reads a header field and returns the value.
func Header(client *clients.Client, key string) string {
	return client.Request.Header.Get(key)
}
