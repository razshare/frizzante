package guard

import (
	"main/lib/core/client"
	"main/lib/core/tag"
)

type Guard struct {
	Name    string
	Handler func(c *client.Client, allow func())
	Tags    []tag.Tag
}
