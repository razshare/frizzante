package route

import "github.com/razshare/frizzante/client"

type Route struct {
	Pattern string
	Handler func(c *client.Client)
	Tags    []string
}
