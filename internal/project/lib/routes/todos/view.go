package todos

import (
	"github.com/razshare/frizzante/internal/project/lib/core/clients"
	"github.com/razshare/frizzante/internal/project/lib/core/receive"
	"github.com/razshare/frizzante/internal/project/lib/core/send"
	"github.com/razshare/frizzante/internal/project/lib/core/views"
	"github.com/razshare/frizzante/internal/project/lib/sessions"
)

func View(client *clients.Client) {
	var session sessions.Session
	if !receive.Session(client, &session) {
		session = *sessions.NewDefault()
	}

	send.View(client, views.View{Name: "todos", Props: Props{
		Error: session.Error,
		Items: session.Todos,
	}})
	session.Error = ""
}
