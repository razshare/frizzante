package guard

import (
	"github.com/razshare/frizzante/internal/project/lib/core/client"
	"github.com/razshare/frizzante/internal/project/lib/core/tag"
)

type Guard struct {
	Name    string
	Handler func(c *client.Client, allow func())
	Tags    []tag.Tag
}
