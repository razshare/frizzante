package receive

import "github.com/razshare/frizzante/internal/project/lib/core/clients"

// QueryExists checks if a query field exists.
func QueryExists(client *clients.Client, key string) bool {
	return client.Request.URL.Query().Has(key)
}
