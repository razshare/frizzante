package fallback

import (
	"github.com/razshare/frizzante/internal/project/lib/core/client"
	"github.com/razshare/frizzante/internal/project/lib/core/send"
	"github.com/razshare/frizzante/internal/project/lib/routes/handlers/welcome"
)

func View(client *client.Client) {
	send.FileOrElse(client, func() { welcome.View(client) })
}
