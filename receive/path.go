package receive

import "github.com/razshare/frizzante/client"

// Path reads a parameters fields and returns the value.
//
// Compatible with web sockets.
func Path(c *client.Client, k string) string {
	return c.Request.PathValue(k)
}
