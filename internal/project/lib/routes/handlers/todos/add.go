package todos

import (
	"github.com/razshare/frizzante/internal/project/lib/core/client"
	"github.com/razshare/frizzante/internal/project/lib/core/receive"
	"github.com/razshare/frizzante/internal/project/lib/core/send"
	"github.com/razshare/frizzante/internal/project/lib/core/view"
	"github.com/razshare/frizzante/internal/project/lib/session/memory"
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
