package guard

import "github.com/razshare/frizzante/conn"

type Guard struct {
	Name    string
	Handler func(c *conn.Conn, allow func())
	Tags    []string
}
