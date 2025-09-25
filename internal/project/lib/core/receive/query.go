package receive

import "github.com/razshare/frizzante/internal/project/lib/core/client"

// Query reads a query field and returns the value.
func Query(client *client.Client, key string) string {
	return client.Request.URL.Query().Get(key)
}
