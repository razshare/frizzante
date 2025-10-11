package todos

import (
	"github.com/razshare/frizzante/internal/project/lib/core/clients"
	"github.com/razshare/frizzante/internal/project/lib/core/receive"
	"github.com/razshare/frizzante/internal/project/lib/core/send"
	"github.com/razshare/frizzante/internal/project/lib/sessions"
)

func Add(client *clients.Client) {
	session := sessions.Start(receive.SessionId(client))

	var description string
	if receive.FormValue(client, "description", &description); description == "" {
		session.Error = "todo description cannot be empty"
		send.Navigate(client, "/todos")
		return
	}

	session.Todos = append(session.Todos, sessions.Todo{
		Checked:     false,
		Description: description,
	})

	send.Navigate(client, "/todos")
}
