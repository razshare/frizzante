package fallback

import (
	_client "github.com/razshare/frizzante/internal/project/lib/core/client"
	"github.com/razshare/frizzante/internal/project/lib/core/send"
	"github.com/razshare/frizzante/internal/project/lib/routes/handlers/welcome"
)

func View(client *_client.Client) {
	send.FileOrElse(client, func() { welcome.View(client) })
}
