package routes

import (
	"github.com/razshare/frizzante/connections"
)

type Route struct {
	Pattern string
	Handler func(connection *connections.Connection)
	Tags    []string
}
