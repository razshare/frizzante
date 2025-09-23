package todos

import (
	"strconv"

	"github.com/razshare/frizzante/internal/project/lib/core/client"
	"github.com/razshare/frizzante/internal/project/lib/core/receive"
	"github.com/razshare/frizzante/internal/project/lib/core/send"
	session "github.com/razshare/frizzante/internal/project/lib/session/memory"
)

func Toggle(client *client.Client) {
	var err error
	var index int64
	var value int64
	var count int64
	var queryIndex string
	var queryValue string

	state := session.Start(receive.SessionId(client))

	if queryIndex = receive.Query(client, "index"); queryIndex == "" {
		// No index found, ignore the request.
		send.Navigate(client, "/todos")
		return
	}

	if queryValue = receive.Query(client, "value"); queryValue == "" {
		// No index found, ignore the request.
		send.Navigate(client, "/todos")
		return
	}

	if index, err = strconv.ParseInt(queryIndex, 10, 64); err != nil {
		send.Navigatef(client, "/todos?error=%s", err.Error())
		return
	}

	if value, err = strconv.ParseInt(queryValue, 10, 64); err != nil {
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
