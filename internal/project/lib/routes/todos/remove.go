package todos

import (
	"github.com/razshare/frizzante/internal/project/lib/core/clients"
	"github.com/razshare/frizzante/internal/project/lib/core/receive"
	"github.com/razshare/frizzante/internal/project/lib/core/send"
	"github.com/razshare/frizzante/internal/project/lib/sessions"
)

func Remove(client *clients.Client) {
	session := sessions.Start(receive.SessionId(client))

	var index int
	if !receive.FormValue(client, "index", &index) {
		session.Error = "could not parse index"
		send.Navigate(client, "/todos")
		return
	}

	if count := len(session.Todos); index >= count || index < 0 {
		session.Error = "index out of bounds"
		send.Navigate(client, "/todos")
		return
	}

	session.Todos = append(
		session.Todos[:index],
		session.Todos[index+1:]...,
	)

	send.Navigate(client, "/todos")
}
