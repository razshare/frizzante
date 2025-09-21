package route

import (
	"github.com/razshare/frizzante/internal/project/lib/core/client"
	"github.com/razshare/frizzante/internal/project/lib/core/tag"
)

type Route struct {
	Pattern string
	Handler func(client *client.Client)
	Tags    []tag.Tag
}
