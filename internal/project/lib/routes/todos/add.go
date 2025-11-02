package todos

import (
	"github.com/razshare/frizzante/internal/project/lib/core/clients"
	"github.com/razshare/frizzante/internal/project/lib/core/receive"
	"github.com/razshare/frizzante/internal/project/lib/core/send"
	"github.com/razshare/frizzante/internal/project/lib/sessions"
)

func Add(client *clients.Client) {
	var session sessions.Session
	if !receive.Session(client, &session) {
		session = *sessions.NewDefault()
	}

	var form FormAdd
	if !receive.Form(client, &form) {
		session.Error = "could not parse form"
		send.Navigate(client, "/todos")
		return
	}

	if form.Description == "" {
		session.Error = "description cannot be empty"
		send.Navigate(client, "/todos")
		return
	}

	session.Todos = append(session.Todos, sessions.Todo{
		Checked:     false,
		Description: form.Description,
	})

	send.Navigate(client, "/todos")
}
