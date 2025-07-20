package routes

import "github.com/razshare/frizzante/connections"

type Route struct {
	Pattern string
	Handler func(c *connections.Connection)
	Tags    []string
}
