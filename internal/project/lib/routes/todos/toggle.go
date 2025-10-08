package todos

import (
	"strconv"

	"github.com/razshare/frizzante/internal/project/lib/core/clients"
	"github.com/razshare/frizzante/internal/project/lib/core/receive"
	"github.com/razshare/frizzante/internal/project/lib/core/send"
	"github.com/razshare/frizzante/internal/project/lib/sessions"
)

func Toggle(client *clients.Client) {
	var err error
	var count int64
	var index int64
	var value int64
	var indexQuery string
	var valueQuery string

	state := sessions.Start(receive.SessionId(client))

	if indexQuery = receive.Query(client, "index"); indexQuery == "" {
		// No index found, ignore the request.
		send.Navigate(client, "/todos")
		return
	}

	if valueQuery = receive.Query(client, "value"); valueQuery == "" {
		// No value found, ignore the request.
		send.Navigate(client, "/todos")
		return
	}

	if index, err = strconv.ParseInt(indexQuery, 10, 64); err != nil {
		// Could not parse index, redirect with error.
		send.Navigatef(client, "/todos?error=%s", err.Error())
		return
	}

	if value, err = strconv.ParseInt(valueQuery, 10, 64); err != nil {
		// Could not parse value, redirect with error.
		send.Navigatef(client, "/todos?error=%s", err.Error())
		return
	}

	if count = int64(len(state.Todos)); index >= count || index < 0 {
		// Index is out of bounds, ignore the request.
		send.Navigate(client, "/todos")
		return
	}

	state.Todos[index].Checked = value > 0

	send.Navigate(client, "/todos")
}
