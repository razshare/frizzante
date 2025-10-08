package todos

import (
	"github.com/razshare/frizzante/internal/project/lib/core/clients"
	"github.com/razshare/frizzante/internal/project/lib/core/receive"
	"github.com/razshare/frizzante/internal/project/lib/core/send"
	"github.com/razshare/frizzante/internal/project/lib/sessions"
)

func Add(client *clients.Client) {
	var description string

	state := sessions.Start(receive.SessionId(client))

	if description = receive.Query(client, "description"); description == "" {
		send.Navigate(client, "/todos?error=todo description cannot be empty")
		return
	}

	state.Todos = append(state.Todos, sessions.Todo{
		Checked:     false,
		Description: description,
	})

	send.Navigate(client, "/todos")
}
