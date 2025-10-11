package todos

import (
	"github.com/razshare/frizzante/internal/project/lib/core/clients"
	"github.com/razshare/frizzante/internal/project/lib/core/receive"
	"github.com/razshare/frizzante/internal/project/lib/core/send"
	"github.com/razshare/frizzante/internal/project/lib/core/views"
	"github.com/razshare/frizzante/internal/project/lib/sessions"
)

func Toggle(client *clients.Client) {
	var value int

	session := sessions.Start(receive.SessionId(client))

	var index int
	if !receive.FormParsedValue(client, "index", &index) {
		send.View(client, views.View{Name: "Todos", Props: Props{
			Error: "could not parse index",
			Items: session.Todos,
		}})
		return
	}

	if !receive.FormParsedValue(client, "value", &value) {
		send.View(client, views.View{Name: "Todos", Props: Props{
			Error: "could not parse value",
			Items: session.Todos,
		}})
		return
	}

	if count := len(session.Todos); index >= count || index < 0 {
		send.View(client, views.View{Name: "Todos", Props: Props{
			Error: "index out of bounds",
			Items: session.Todos,
		}})
		return
	}

	session.Todos[index].Checked = value > 0

	send.View(client, views.View{Name: "Todos", Props: Props{
		Items: session.Todos,
	}})
}
