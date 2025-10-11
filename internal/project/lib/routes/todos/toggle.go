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
	if !receive.FormValue(client, "index", &index) {
		session.Error = "could not parse index"
		send.Navigate(client, "/todos")
		return
	}

	var value int
	if !receive.FormValue(client, "value", &value) {
		session.Error = "could not parse value"
		send.Navigate(client, "/todos")
		return
	}

	if count := len(session.Todos); index >= count || index < 0 {
		session.Error = "index out of bounds"
		send.Navigate(client, "/todos")
		return
	}

	session.Todos[index].Checked = value > 0

	send.Navigate(client, "/todos")
}
