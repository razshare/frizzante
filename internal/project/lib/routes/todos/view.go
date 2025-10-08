package todos

import (
	"github.com/razshare/frizzante/internal/project/lib/core/clients"
	"github.com/razshare/frizzante/internal/project/lib/core/receive"
	"github.com/razshare/frizzante/internal/project/lib/core/send"
	"github.com/razshare/frizzante/internal/project/lib/core/views"
	"github.com/razshare/frizzante/internal/project/lib/sessions"
)

func View(client *clients.Client) {
	session := sessions.Start(receive.SessionId(client))
	send.View(client, views.View{
		Name: "Todos",
		Props: Props{
			Items: session.Todos,
			Error: receive.Query(client, "error"),
		},
	})
}
