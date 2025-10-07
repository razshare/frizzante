package todos

import (
	_client "github.com/razshare/frizzante/internal/project/lib/core/client"
	"github.com/razshare/frizzante/internal/project/lib/core/receive"
	"github.com/razshare/frizzante/internal/project/lib/core/send"
	_session "github.com/razshare/frizzante/internal/project/lib/session"
)

func Add(client *_client.Client) {
	var description string

	state := _session.Start(receive.SessionId(client))

	if description = receive.Query(client, "description"); description == "" {
		send.Navigate(client, "/todos?error=todo description cannot be empty")
		return
	}

	state.Todos = append(state.Todos, _session.Todo{
		Checked:     false,
		Description: description,
	})

	send.Navigate(client, "/todos")
}
