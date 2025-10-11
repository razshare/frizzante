package todos

import (
	"github.com/razshare/frizzante/internal/project/lib/core/clients"
	"github.com/razshare/frizzante/internal/project/lib/core/receive"
	"github.com/razshare/frizzante/internal/project/lib/core/send"
	"github.com/razshare/frizzante/internal/project/lib/core/views"
	"github.com/razshare/frizzante/internal/project/lib/sessions"
)

func Add(client *clients.Client) {
	var description string

	session := sessions.Start(receive.SessionId(client))

	if ok := receive.FormParsedValue(client, "description", &description); !ok || description == "" {
		send.View(client, views.View{Name: "Todos", Props: Props{
			Error: "todo description cannot be empty",
			Items: session.Todos,
		}})
		return
	}

	session.Todos = append(session.Todos, sessions.Todo{
		Checked:     false,
		Description: description,
	})

	send.View(client, views.View{Name: "Todos", Props: Props{
		Items: session.Todos,
	}})
}
