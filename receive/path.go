package receive

import "github.com/razshare/frizzante/client"

// Path reads a parameters fields and returns the value.
//
// Compatible with web sockets.
func Path(client *client.Client, key string) string {
	return client.Request.PathValue(key)
}
