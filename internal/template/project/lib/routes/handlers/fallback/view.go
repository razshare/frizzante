package fallback

import (
	"main/lib/routes/handlers/welcome"

	"github.com/razshare/frizzante/client"
	"github.com/razshare/frizzante/send"
)

func View(c *client.Client) {
	send.FileOrElse(c, func() { welcome.View(c) })
}
