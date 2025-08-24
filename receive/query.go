package receive

import "github.com/razshare/frizzante/client"

// Query reads a query field and returns the value.
//
// Compatible with web sockets.
func Query(c *client.Client, k string) string {
	return c.Request.URL.Query().Get(k)
}
