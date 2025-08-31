package route

import (
	"main/lib/core/client"
	"main/lib/core/tag"
)

type Route struct {
	Pattern string
	Handler func(c *client.Client)
	Tags    []tag.Tag
}
