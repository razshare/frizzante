package fallback

import (
	"github.com/razshare/frizzante/internal/project/lib/core/client"
	"github.com/razshare/frizzante/internal/project/lib/core/send"
	"github.com/razshare/frizzante/internal/project/lib/routes/handlers/welcome"
)

// var Dev = os.Getenv("DEV") == "1"

func View(client *client.Client) {
	send.FileOrElse(client, true, func() { welcome.View(client) })
}
