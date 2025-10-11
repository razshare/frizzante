package todos

import (
	"github.com/razshare/frizzante/internal/project/lib/core/clients"
	"github.com/razshare/frizzante/internal/project/lib/core/receive"
	"github.com/razshare/frizzante/internal/project/lib/core/send"
	"github.com/razshare/frizzante/internal/project/lib/core/views"
	"github.com/razshare/frizzante/internal/project/lib/sessions"
)

func Remove(client *clients.Client) {
	var index int

	session := sessions.Start(receive.SessionId(client))

	if !receive.FormParsedValue(client, "index", &index) {
		send.View(client, views.View{Name: "Todos", Props: Props{
			Error: "could not parse index",
			Items: session.Todos,
		}})
		return
	}

	if count := len(session.Todos); index >= count || index < 0 {
		send.View(client, views.View{Name: "Todos", Props: Props{
			Error: "index out of bounds",
			Items: session.Todos,
		}})
		return
	}

	session.Todos = append(
		session.Todos[:index],
		session.Todos[index+1:]...,
	)

	send.View(client, views.View{Name: "Todos", Props: Props{
		Items: session.Todos,
	}})
}
