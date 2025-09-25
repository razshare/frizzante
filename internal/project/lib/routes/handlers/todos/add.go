package todos

import (
	"github.com/razshare/frizzante/internal/project/lib/core/client"
	"github.com/razshare/frizzante/internal/project/lib/core/receive"
	"github.com/razshare/frizzante/internal/project/lib/core/send"
	session "github.com/razshare/frizzante/internal/project/lib/session/memory"
)

func Add(client *client.Client) {
	var description string

	state := session.Start(receive.SessionId(client))

	if description = receive.Query(client, "description"); description == "" {
		send.Navigate(client, "/todos?error=todo description cannot be empty")
		return
	}

	state.Todos = append(state.Todos, session.Todo{
		Checked:     false,
		Description: description,
	})

	send.Navigate(client, "/todos")
}
