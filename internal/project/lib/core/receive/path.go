package receive

import "github.com/razshare/frizzante/internal/project/lib/core/clients"

// Path reads a path value field and returns it.
func Path(client *clients.Client, key string) string {
	return client.Request.PathValue(key)
}
