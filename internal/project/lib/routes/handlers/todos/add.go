package todos

import (
	"github.com/razshare/frizzante/internal/project/lib/core/client"
	"github.com/razshare/frizzante/internal/project/lib/core/receive"
	"github.com/razshare/frizzante/internal/project/lib/core/send"
	"github.com/razshare/frizzante/internal/project/lib/session/memory"
)

func Add(client *client.Client) {
	var query string
	var state *session.State

	state = session.Start(receive.SessionId(client))

	if query = receive.Query(client, "description"); query == "" {
		send.Navigate(client, "/todos?error=todo description cannot be empty")
		return
	}

	state.Todos = append(state.Todos, session.Todo{
		Checked:     false,
		Description: query,
	})

	send.Navigate(client, "/todos")
}
