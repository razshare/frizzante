package todos

import (
	"strconv"

	"github.com/razshare/frizzante/internal/project/lib/core/client"
	"github.com/razshare/frizzante/internal/project/lib/core/receive"
	"github.com/razshare/frizzante/internal/project/lib/core/send"
	"github.com/razshare/frizzante/internal/project/lib/session/memory"
)

func Remove(client *client.Client) {
	state := session.Start(receive.SessionId(client))

	count := int64(len(state.Todos))
	if count == 0 {
		// No index found, ignore the request.
		send.Navigate(client, "/todos")
		return
	}

	query := receive.Query(client, "index")
	if query == "" {
		// No index found, ignore the request.
		send.Navigate(client, "/todos")
		return
	}

	index, err := strconv.ParseInt(query, 10, 64)
	if err != nil {
		send.Navigatef(client, "/todos?error=%s", err.Error())
		return
	}
	if index >= count {
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
