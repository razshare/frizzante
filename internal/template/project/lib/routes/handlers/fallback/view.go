package fallback

import (
	"github.com/razshare/frizzante/client"
	"github.com/razshare/frizzante/send"
	"main/lib/routes/handlers/welcome"
)

func View(c *client.Client) {
	send.FileOrElse(c, func() { welcome.View(c) })
}
