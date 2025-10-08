package todos

import (
	"github.com/razshare/frizzante/internal/project/lib/core/clients"
	"github.com/razshare/frizzante/internal/project/lib/core/receive"
	"github.com/razshare/frizzante/internal/project/lib/core/send"
	"github.com/razshare/frizzante/internal/project/lib/sessions"
)

func Remove(client *clients.Client) {
	state := sessions.Start(receive.SessionId(client))

	var index int
	if !receive.FormInt(client, "index", &index) {
		// Could not parse index, redirect with error.
		send.Navigate(client, "/todos?error=could not parse index")
	}

	if count := len(state.Todos); index >= count || index < 0 {
		// Index is out of bounds, ignore the request.
		send.Navigate(client, "/todos")
		return
	}

	state.Todos = append(
		state.Todos[:index],
		state.Todos[index+1:]...,
	)

	send.Navigate(client, "/todos")
}
