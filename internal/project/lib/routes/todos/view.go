package todos

import (
	"github.com/razshare/frizzante/internal/project/lib/core/clients"
	"github.com/razshare/frizzante/internal/project/lib/core/send"
	"github.com/razshare/frizzante/internal/project/lib/core/views"
	"github.com/razshare/frizzante/internal/project/lib/memory/sessions"
)

func View(client *clients.Client) {
	session := sessions.Start(client)
	defer func() { session.Error = "" }()
	send.View(client, views.View{Name: "Todos", Props: Props{
		Error: session.Error,
		Items: session.Todos,
	}})
}
