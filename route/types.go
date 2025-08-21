package route

import (
	"github.com/razshare/frizzante/client"
	"github.com/razshare/frizzante/tag"
)

type Route struct {
	Pattern string
	Handler func(c *client.Client)
	Tags    []tag.Tag
}
