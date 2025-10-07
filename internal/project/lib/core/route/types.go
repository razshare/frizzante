package route

import (
	_client "github.com/razshare/frizzante/internal/project/lib/core/client"
	"github.com/razshare/frizzante/internal/project/lib/core/tag"
)

type Route struct {
	Pattern string
	Handler func(client *_client.Client)
	Tags    []tag.Tag
}
