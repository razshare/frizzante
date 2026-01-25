package todos

import (
	"github.com/razshare/frizzante/internal/project/lib/core/clients"
	"github.com/razshare/frizzante/internal/project/lib/core/receive"
	"github.com/razshare/frizzante/internal/project/lib/core/send"
	"github.com/razshare/frizzante/internal/project/lib/sessions"
)

func Remove(client *clients.Client) {
	session := sessions.NewDefault()
	receive.Session(client, &session)
	var form FormRemove
	if !receive.Form(client, &form) {
		session.Error = "could not parse form"
		send.Navigate(client, "/todos")
		return
	}
	if count := len(session.Todos); form.Index >= count || form.Index < 0 {
		session.Error = "index out of bounds"
		send.Navigate(client, "/todos")
		return
	}
	session.Todos = append(
		session.Todos[:form.Index],
		session.Todos[form.Index+1:]...,
	)
	send.Navigate(client, "/todos")
}
