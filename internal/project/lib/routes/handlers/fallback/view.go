package fallback

import (
	"os"

	"github.com/razshare/frizzante/internal/project/lib/core/client"
	"github.com/razshare/frizzante/internal/project/lib/core/send"
	"github.com/razshare/frizzante/internal/project/lib/routes/handlers/welcome"
)

var dev = os.Getenv("DEV") == "1"

func View(client *client.Client) {
	send.FileOrElse(client, dev, func() { welcome.View(client) })
}
