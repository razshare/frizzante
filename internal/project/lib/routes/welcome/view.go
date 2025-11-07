package welcome

import (
	"github.com/razshare/frizzante/internal/project/lib/core/clients"
	"github.com/razshare/frizzante/internal/project/lib/core/send"
	"github.com/razshare/frizzante/internal/project/lib/core/views"
)

func View(client *clients.Client) {
	send.View(client, views.View{Name: "Welcome"})
}
