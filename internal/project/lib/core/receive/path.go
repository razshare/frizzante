package receive

import _client "github.com/razshare/frizzante/internal/project/lib/core/client"

// Path reads a parameters fields and returns the value.
func Path(client *_client.Client, key string) string {
	return client.Request.PathValue(key)
}
