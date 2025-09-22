package receive

import (
	"net/url"

	"github.com/razshare/frizzante/internal/project/lib/core/client"
)

type MultipartForm struct {
	url.Values
	Client *client.Client
}
