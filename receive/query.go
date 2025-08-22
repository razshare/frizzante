package receive

import "github.com/razshare/frizzante/client"

// Query reads a query field and returns the value.
//
// Compatible with web sockets.
func Query(client *client.Client, k string) string {
	return client.Request.URL.Query().Get(k)
}
