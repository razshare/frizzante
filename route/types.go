package route

import "github.com/razshare/frizzante/conn"

type Route struct {
	Pattern string
	Handler func(c *conn.Conn)
	Tags    []string
}
