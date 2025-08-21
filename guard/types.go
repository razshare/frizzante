package guard

import (
	"github.com/razshare/frizzante/client"
	"github.com/razshare/frizzante/tag"
)

type Guard struct {
	Name    string
	Handler func(c *client.Client, allow func())
	Tags    []tag.Tag
}
