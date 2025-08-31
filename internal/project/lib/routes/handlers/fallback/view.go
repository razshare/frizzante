package fallback

import (
	"main/lib/routes/handlers/welcome"

	"main/lib/core/client"
	"main/lib/core/send"
)

func View(c *client.Client) {
	send.FileOrElse(c, func() { welcome.View(c) })
}
