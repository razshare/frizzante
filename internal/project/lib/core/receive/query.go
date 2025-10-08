package receive

import "github.com/razshare/frizzante/internal/project/lib/core/clients"

// Query reads a query field and returns the value.
func Query(client *clients.Client, key string) string {
	return client.Request.URL.Query().Get(key)
}
