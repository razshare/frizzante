package receive

import "github.com/razshare/frizzante/conn"

// Query reads a query field and returns the value.
//
// Compatible with web sockets.
func Query(c *conn.Conn, k string) string {
	return c.Request.URL.Query().Get(k)
}
