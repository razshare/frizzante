package guard

import (
	_client "github.com/razshare/frizzante/internal/project/lib/core/client"
	"github.com/razshare/frizzante/internal/project/lib/core/tag"
)

type Guard struct {
	Name    string
	Handler func(client *_client.Client, allow func())
	Tags    []tag.Tag
}
