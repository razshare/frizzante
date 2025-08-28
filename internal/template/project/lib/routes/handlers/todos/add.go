package todos

import (
	"github.com/razshare/frizzante/client"
	"github.com/razshare/frizzante/receive"
	"github.com/razshare/frizzante/send"
	"github.com/razshare/frizzante/view"
	"main/lib/session"
)

func Add(c *client.Client) {
	s := session.Start(receive.SessionId(c))
	d := receive.Query(c, "description")
	if d == "" {
		send.View(c, view.View{
			Name: "Todos",
			Props: map[string]any{
				"todos": s.Todos,
				"error": "todo description cannot be empty",
			},
		})
		return
	}

	s.Todos = append(s.Todos, session.Todo{
		Checked:     false,
		Description: d,
	})

	send.Navigate(c, "/todos")
}
