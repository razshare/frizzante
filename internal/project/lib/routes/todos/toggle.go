package todos

import (
	"github.com/razshare/frizzante/internal/project/lib/core/clients"
	"github.com/razshare/frizzante/internal/project/lib/core/receive"
	"github.com/razshare/frizzante/internal/project/lib/core/send"
	"github.com/razshare/frizzante/internal/project/lib/sessions"
)

func Toggle(client *clients.Client) {
	session := sessions.Start(receive.SessionId(client))

	var index int
	if !receive.FormInt(client, "index", &index) {
		send.Navigate(client, "/todos?error=could not parse index")
		return
	}

	var value int
	if !receive.FormInt(client, "value", &value) {
		send.Navigate(client, "/todos?error=could not parse value")
		return
	}

	if count := len(session.Todos); index >= count || index < 0 {
		// Index is out of bounds, ignore the request.
		send.Navigate(client, "/todos")
		return
	}

	session.Todos[index].Checked = value > 0

	send.Navigate(client, "/todos")
}
