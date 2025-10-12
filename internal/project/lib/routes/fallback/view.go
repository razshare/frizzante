package fallback

import (
	"github.com/razshare/frizzante/internal/project/lib/core/clients"
	"github.com/razshare/frizzante/internal/project/lib/core/send"
	"github.com/razshare/frizzante/internal/project/lib/routes/welcome"
)

func View(client *clients.Client) {
	if !send.RequestedFile(client) {
		welcome.View(client)
	}
}
