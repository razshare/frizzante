package fallback

import (
	"os"

	"github.com/razshare/frizzante/internal/project/lib/core/client"
	"github.com/razshare/frizzante/internal/project/lib/core/send"
	"github.com/razshare/frizzante/internal/project/lib/routes/handlers/welcome"
)

func View(client *client.Client) {
	send.FileOrElse(client, send.FileOrElseConfig{
		UseDisk: os.Getenv("DEV") == "1",
		OrElse:  func() { welcome.View(client) },
	})
}
