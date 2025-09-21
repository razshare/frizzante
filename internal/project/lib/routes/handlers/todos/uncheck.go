package todos

import (
	"strconv"

	"github.com/razshare/frizzante/internal/project/lib/core/client"
	"github.com/razshare/frizzante/internal/project/lib/core/receive"
	"github.com/razshare/frizzante/internal/project/lib/core/send"
	"github.com/razshare/frizzante/internal/project/lib/session/memory"
)

func Uncheck(client *client.Client) {
	state := session.Start(receive.SessionId(client))

	query := receive.Query(client, "index")
	if query == "" {
		// No index found, ignore the request.
		send.Navigate(client, "/todos")
		return
	}

	index, err := strconv.ParseInt(query, 10, 64)
	if nil != err {
		send.Navigatef(client, "/todos?error=%s", err.Error())
		return
	}

	count := int64(len(state.Todos))

	if index >= count {
		// Index is out of bounds, ignore the request.
		send.Navigate(client, "/todos")
		return
	}

	state.Todos[index].Checked = false

	send.Navigate(client, "/todos")
}
